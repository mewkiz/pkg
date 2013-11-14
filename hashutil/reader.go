package hashutil

import (
	"io"
)

// Hash16Reader is the interface that wraps the Hash16 and io.Reader interfaces.
// All data read from the io.Reader is also added to the running hash, unless
// Hash16 is set to nil.
type Hash16Reader interface {
	Hash16
	io.Reader
}

// hash16Reader is a simple implementation of the Hash16Reader interface.
type hash16Reader struct {
	Hash16
	io.Reader
}

// NewHash16Reader returns a Hash16Reader based on the provided io.Reader. All
// data read is also added to the running hash, unless Hash16 is set to nil.
func NewHash16Reader(r io.Reader, h Hash16) (hr Hash16Reader) {
	hr = &hash16Reader{
		Hash16: h,
		Reader: r,
	}
	return hr
}

// Read reads up to len(buf) bytes into buf. All data read is also added to the
// running hash, unless Hash16 is set to nil. It returns the number of bytes
// read (0 <= n <= len(buf)) and any error encountered.
func (hr *hash16Reader) Read(buf []byte) (n int, err error) {
	if hr.Hash16 == nil {
		return hr.Reader.Read(buf)
	}
	n, err = hr.Reader.Read(buf)
	if n > 0 {
		_, e := hr.Hash16.Write(buf[:n])
		if err == nil && e != nil {
			return n, e
		}
	}
	return n, err
}
