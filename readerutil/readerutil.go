// Package readerutil implements io.Reader utility functions.
package readerutil

import "io"
import "os"

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

// Peek returns the next n bytes without advancing the reader. If Peek returns
// fewer than n bytes, it also returns io.ErrUnexpectedEOF.
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
		// Read error.
		return nil, err
	}

	if m < n {
		// Short read.
		buf = buf[:m]
	}

	return buf, e
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
