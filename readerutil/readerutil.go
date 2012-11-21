// Package readerutil implements io.Reader utility functions.
package readerutil

import "encoding/binary"
import "io"
import "os"
import "unicode"
import "unicode/utf16"
import "unicode/utf8"

// NewByteReader returns a new io.ByteReader based on the provided io.Reader.
func NewByteReader(r io.Reader) io.ByteReader {
	return byteReader{Reader: r}
}

// byteReader implements the io.ByteReader interface.
type byteReader struct {
	io.Reader
}

// ReadByte reads and returns the next byte from the input. If no byte is
// available, err will be set.
func (r byteReader) ReadByte() (c byte, err error) {
	return ReadByte(r)
}

// ReadByte reads and returns the next byte from the provided io.Reader.
func ReadByte(r io.Reader) (c byte, err error) {
	buf := make([]byte, 1)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return 0, err
	}

	return buf[0], nil
}

// A Peeker is an io.ReadSeeker that can also peek ahead.
//
// Peek returns the next n bytes without advancing the reader. The error is EOF
// only if no bytes were read. If an EOF happens after reading some but not all
// the bytes, ReadFull returns ErrUnexpectedEOF.
type Peeker interface {
	io.ReadSeeker
	Peek(n int) (buf []byte, err error)
}

// peeker implements the Peeker interface.
type peeker struct {
	io.ReadSeeker
}

// NewPeeker returns a new Peeker based on the provided io.ReadSeeker.
func NewPeeker(r io.ReadSeeker) Peeker {
	return peeker{ReadSeeker: r}
}

// Peek returns the next n bytes without advancing the reader. The error is EOF
// only if no bytes were read. If an EOF happens after reading some but not all
// the bytes, ReadFull returns ErrUnexpectedEOF.
func (r peeker) Peek(n int) (buf []byte, err error) {
	// Record original position.
	orig, err := r.Seek(0, os.SEEK_CUR)
	if err != nil {
		return nil, err
	}

	// Read content, but check error after position reset.
	buf = make([]byte, n)
	m, e := io.ReadFull(r, buf)

	// Reset original position.
	_, err = r.Seek(orig, os.SEEK_SET)
	if err != nil {
		return nil, err
	}

	if e != nil && e != io.ErrUnexpectedEOF {
		// Return read error.
		return nil, e
	}

	if m < n {
		// Short read.
		buf = buf[:m]
	}

	return buf, e
}

// A BinaryPeeker is an io.ReadSeeker that can also peek ahead.
type BinaryPeeker interface {
	io.ReadSeeker
	Peek(data interface{}) (err error)
}

// binaryPeeker implements the BinaryPeeker interface.
type binaryPeeker struct {
	io.ReadSeeker
	order binary.ByteOrder
}

// NewBinaryPeeker returns a new BinaryPeeker based on the provided
// io.ReadSeeker.
//
// Peek reads structured binary data without advancing the reader. Data must be
// a pointer to a fixed-size value or a slice of fixed-size values. Bytes read
// from r are decoded using the receiver's byte order and written to successive
// fields of the data. When reading into structs, the field data for fields with
// blank (_) field names is skipped; i.e., blank field names may be used for
// padding.
func NewBinaryPeeker(r io.ReadSeeker, order binary.ByteOrder) BinaryPeeker {
	return binaryPeeker{ReadSeeker: r, order: order}
}

// Peek reads structured binary data without advancing the reader. Data must be
// a pointer to a fixed-size value or a slice of fixed-size values. Bytes read
// from r are decoded using the receiver's byte order and written to successive
// fields of the data. When reading into structs, the field data for fields with
// blank (_) field names is skipped; i.e., blank field names may be used for
// padding.
func (r binaryPeeker) Peek(data interface{}) (err error) {
	// Record original position.
	orig, err := r.Seek(0, os.SEEK_CUR)
	if err != nil {
		return err
	}

	// Read content, but check error after position reset.
	e := binary.Read(r, r.order, data)

	// Reset original position.
	_, err = r.Seek(orig, os.SEEK_SET)
	if err != nil {
		return err
	}

	if e != nil {
		// Return read error.
		return e
	}

	return nil
}

// Size returns the total size in bytes of the provided io.Seeker. The original
// position is preserved.
func Size(r io.Seeker) (n int64, err error) {
	// Record original position.
	orig, err := r.Seek(0, os.SEEK_CUR)
	if err != nil {
		return 0, err
	}

	// Seek end position.
	end, err := r.Seek(0, os.SEEK_END)
	if err != nil {
		return 0, err
	}

	// Reset original position.
	_, err = r.Seek(orig, os.SEEK_SET)
	if err != nil {
		return 0, err
	}

	return end, nil
}

// IsUTF8 decodes a chunk of data as UTF-8 and returns true if no decoding
// errors occured and at least 75% of the decoded runes are graphic or space and
// within the ASCII range.
func IsUTF8(r io.ReadSeeker) (ok bool, err error) {
	// Read a chunk of at most 256 bytes.
	const chunkSize = 256
	rr := NewPeeker(r)
	buf, err := rr.Peek(chunkSize)
	if err != nil && err != io.ErrUnexpectedEOF {
		return false, err
	}

	// Decode the chunk as UTF-8. Make sure that no short rune reads are
	// performed.
	var printableCount, total float64
	for utf8.FullRune(buf) {
		r, size := utf8.DecodeRune(buf)
		if r == utf8.RuneError {
			// No rune errors should occur in valid UTF-8.
			return false, nil
		}
		if isPrintable(r) {
			// Record printable runes.
			printableCount++
		}
		total++
		buf = buf[size:]
	}

	if total < 1 {
		// No valid UTF-8 runes located.
		return false, nil
	}

	if printableCount/total >= 0.75 {
		// Assume that the data is UTF-8, since at least 75% of the runes are
		// graphic or space and within the ASCII range.
		return true, nil
	}

	return false, nil
}

// IsUTF16 decodes a chunk of data as UTF-16 with the specified byte order and
// returns true if a valid BOM byte sequence was located or no decoding errors
// occured and at least 75% of the decoded runes are graphic or space and within
// the ASCII range.
func IsUTF16(r io.ReadSeeker, order binary.ByteOrder) (ok bool, err error) {
	// Verify file size.
	size, err := Size(r)
	if err != nil {
		return false, err
	}
	if size < 2 || size%2 != 0 {
		// UTF-16 must be dividable by 2.
		return false, nil
	}

	// Peek for UTF-16 BOM byte sequence.
	const BOM = 0xFEFF
	var bom uint16
	rr := NewBinaryPeeker(r, order)
	err = rr.Peek(&bom)
	if err != nil {
		return false, err
	}
	if bom == BOM {
		return true, nil
	}

	// Read a chunk of at most 256 bytes.
	chunkSize := 256
	if size < int64(chunkSize) {
		chunkSize = int(size)
	}
	buf := make([]uint16, chunkSize/2)
	err = rr.Peek(buf)
	if err != nil {
		return false, err
	}

	// Make sure that no short rune reads are performed.
	const (
		// 0xd800-0xdc00 encodes the high 10 bits of a pair.
		// 0xdc00-0xe000 encodes the low 10 bits of a pair.
		// the value is those 20 bits plus 0x10000.
		surr1 = 0xd800
		surr2 = 0xdc00
		surr3 = 0xe000
	)
	if surr1 <= buf[len(buf)-1] && buf[len(buf)-1] < surr2 {
		// Ignore the last rune if it is the first surrogate of a pair.
		buf = buf[:len(buf)-1]
	}

	// Decode the chunk as UTF-16.
	var printableCount, total float64
	for _, r := range utf16.Decode(buf) {
		if r == unicode.ReplacementChar {
			// No rune errors should occur in valid UTF-16.
			return false, nil
		}
		if isPrintable(r) {
			// Record printable runes.
			printableCount++
		}
		total++
	}

	if total < 1 {
		// No valid UTF-8 runes located.
		return false, nil
	}

	if printableCount/total >= 0.75 {
		// Assume that the data is UTF-16, since at least 75% of the runes are
		// graphic or space and within the ASCII range.
		return true, nil
	}

	return false, nil
}

// isPrintable returns true if the provided rune is graphic or space and within
// the ASCII range.
func isPrintable(r rune) bool {
	if r > unicode.MaxASCII {
		return false
	}
	return unicode.IsGraphic(r) || unicode.IsSpace(r)
}
