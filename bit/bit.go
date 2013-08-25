// Package bit implements convenient handling of bitstreams.
package bit

import (
	"errors"
	"fmt"
	"io"
)

// Stream is a sequence of bits, a bitstream.
type Stream []int8

// NewStream returns a new bitstream based on the byte slice's bits.
func NewStream(buf []byte) (bits Stream) {
	// Allocate a new bitstream.
	n := len(buf) * 8
	bits = make(Stream, 0, n)

	// Append the byte slice's bits.
	bits.appendBytes(buf)

	return bits
}

// AppendBytes appends the byte slice's bits to the end of the bitstream and
// returns the updated bitstream. It is therefore necessary to store the result
// of AppendBytes, often in the variable holding the bitstream itself:
//    bits = bits.AppendBytes(buf)
func (bits Stream) AppendBytes(buf []byte) Stream {
	// Allocate a new bitstream.
	n := len(bits) + len(buf)*8
	dst := make(Stream, len(bits), n)

	// Copy old bits.
	copy(dst, bits)

	// Append the byte slice's bits.
	dst.appendBytes(buf)

	return dst
}

// appendBytes appends the byte slice's bits to the end of the bitstream.
func (dst *Stream) appendBytes(buf []byte) {
	// Append the byte slice's bits.
	for _, b := range buf {
		for i := 0; i < 8; i++ {
			var bit int8
			if b&0x80 != 0 {
				bit = 1
			}
			*dst = append(*dst, bit)
			b <<= 1
		}
	}
}

// Uint64 returns the uint64 representation of bits. It panics if there are more
// than 64 bits.
func (bits Stream) Uint64() (x uint64) {
	if len(bits) > 64 {
		err := fmt.Errorf("Stream.Uint64: too many bits (%d).", len(bits))
		panic(err.Error())
	}
	for _, bit := range bits {
		x <<= 1
		x += uint64(bit)
	}
	return x
}

// Stream satisfies the fmt.Stringer interface.
func (bits Stream) String() string {
	// Reserve space for padding.
	padCount := (len(bits) - 1) / 8
	buf := make([]byte, len(bits)+padCount)

	// Add padding after every 8th bit.
	var pad int
	for i, bit := range bits {
		buf[i+pad] = byte(bit + '0')
		if (i+1)%8 == 0 && i != len(bits)-1 {
			pad++
			buf[i+pad] = ' '
		}
	}

	return string(buf)
}

// Equal returns a boolean reporting whether a == b. A nil argument is
// equivalent to an empty bitstream.
func Equal(a, b Stream) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Reader is the interface that wraps the Read method.
//
// Read reads exactly n bits. It returns the bits as a bitstream and an error if
// fewer bits were read. The error is io.EOF only if no bits were read. If an
// EOF happens after reading some but not all the bits, Read returns
// io.ErrUnexpectedEOF.
type Reader interface {
	Read(n int) (bits Stream, err error)
}

// ReadSeeker is the interface that groups the Read and Seek methods.
type ReadSeeker interface {
	Reader
	io.Seeker
}

// reader satisfies the ReadSeeker interface.
type reader struct {
	r io.ReadSeeker
	// rest contains left over bits from the previous read.
	rest Stream
}

// NewReader returns a bit.ReadSeeker based on the provided io.ReadSeeker.
func NewReader(r io.ReadSeeker) (br ReadSeeker) {
	br = &reader{
		r: r,
	}
	return br
}

// Read reads exactly n bits. It returns the bits as a bitstream and an error if
// fewer bits were read. The error is io.EOF only if no bits were read. If an
// EOF happens after reading some but not all the bits, Read returns
// io.ErrUnexpectedEOF.
func (br *reader) Read(n int) (bits Stream, err error) {
	// Use bits from previous read.
	if n <= len(br.rest) {
		// Return n bits from previous read.
		bits = br.rest[:n]
		br.rest = br.rest[n:]
		return bits, nil
	} else if len(br.rest) > 0 {
		// Add all bits from previous read.
		bits = br.rest
		n -= len(br.rest)
		br.rest = nil
	}

	// Read enough bytes for n bits.
	byteCount := n / 8
	if n%8 != 0 {
		// One more byte needed for extra bits.
		byteCount++
	}
	buf := make([]byte, byteCount)
	m, err := io.ReadFull(br.r, buf)
	if err != nil {
		var part Stream
		if m > 0 {
			// Get bits from the partial read.
			part = NewStream(buf[:m])
		}
		if len(bits) > 0 {
			// Append bits from partial read to bits from previous read.
			bits = append(bits, part...)
			return bits, io.ErrUnexpectedEOF
		}
		return part, err
	}

	// Read bits based on byteCount.
	a := NewStream(buf)
	if len(bits) > 0 {
		// Append bits from read to bits from previous read.
		bits = append(bits, a[:n]...)
	} else {
		bits = a[:n]
	}
	br.rest = a[n:]

	return bits, nil
}

// Seek whence values.
const (
	SeekSet = 0 // Seek relative to the origin of the file.
	SeekCur = 1 // Seek relative to the current offset.
	SeekEnd = 2 // Seek relative to the end.
)

// Seek sets the bit offset for the next Read to bitOff, interpreted according
// to whence: 0 means relative to the origin of the file, 1 means relative to
// the current bit offset, and 2 means relative to the end. It returns the new
// bit offset and an error, if any.
func (br *reader) Seek(bitOff int64, whence int) (bitRet int64, err error) {
	// Calculate absolute bit offset.
	absBitOff, err := br.getAbsOffset(bitOff, whence)
	if absBitOff < 0 {
		return 0, fmt.Errorf("Seek: negative offset (%d).", absBitOff)
	}

	// Seek absolute bit offset.
	br.rest = nil
	byteCount := absBitOff / 8
	bitCount := absBitOff % 8
	byteRet, err := br.r.Seek(byteCount, SeekSet)
	if err != nil {
		return 0, err
	}
	if bitCount != 0 {
		// If bitCount is non-zero, seek that many bits forward. This is done as
		// follows:
		//    - Read 8 bits.
		//    - Set br.rest to those bits, but ignore the first bitCount bits.
		bits, err := br.Read(8)
		if err != nil {
			byteEnd, err := br.r.Seek(0, SeekEnd)
			if err != nil {
				return 0, err
			}
			// calculate absolute end bit offset.
			bitEnd := byteEnd * 8
			return 0, fmt.Errorf("Seek: offset out of range; max %d, got %d.", bitEnd, absBitOff)
		}
		br.rest = bits[bitCount:]
	}

	// Validate new bit offset.
	bitRet = byteRet*8 + bitCount
	if bitRet != absBitOff {
		return 0, fmt.Errorf("Seek: inaccurate offset after seek; expected %d, got %d.", absBitOff, bitRet)
	}

	return bitRet, nil
}

// getAbsOffset calculates and returns the absolute bit offset based on the
// provided relative bit offset and whence.
func (br *reader) getAbsOffset(bitOff int64, whence int) (bitRet int64, err error) {
	switch whence {
	case SeekSet:
		// bitOff is already an absolute bit offset.
		return bitOff, nil

	case SeekCur:
		byteCur, err := br.r.Seek(0, SeekCur)
		if err != nil {
			return 0, err
		}
		// calculate current absolute bit offset.
		bitCur := byteCur*8 - int64(len(br.rest))
		// add relative bit offset.
		bitCur += bitOff
		return bitCur, nil

	case SeekEnd:
		byteEnd, err := br.r.Seek(0, SeekEnd)
		if err != nil {
			return 0, err
		}
		// calculate absolute end bit offset.
		bitEnd := byteEnd * 8
		// add relative bit offset.
		bitEnd += bitOff
		return bitEnd, nil
	}

	return 0, errors.New("Seek: invalid whence.")
}
