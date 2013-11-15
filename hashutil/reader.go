package hashutil

import (
	"hash"
	"io"
)

// hashReader wraps the basic hash.Hash and io.Reader interfaces. All data read
// from the io.Reader is also added to the running hash, unless Hash is set to
// nil.
type hashReader struct {
	hash.Hash
	io.Reader
}

// NewHashReader returns a HashReader based on the provided io.Reader. All data
// read is also added to the running hash, unless Hash is set to nil.
func NewHashReader(r io.Reader, h hash.Hash) (hr *hashReader) {
	hr = &hashReader{
		Hash:   h,
		Reader: r,
	}
	return hr
}

// Read reads up to len(buf) bytes into buf. All data read is also added to the
// running hash, unless Hash is set to nil. It returns the number of bytes read
// (0 <= n <= len(buf)) and any error encountered.
func (hr *hashReader) Read(buf []byte) (n int, err error) {
	if hr.Hash == nil {
		return hr.Reader.Read(buf)
	}
	n, err = hr.Reader.Read(buf)
	if n > 0 {
		_, e := hr.Hash.Write(buf[:n])
		if err == nil && e != nil {
			return n, e
		}
	}
	return n, err
}

// Sum8 returns the 8-bit checksum of the hash. It panics if hr.Hash doesn't
// implement the interface Hash8.
func (hr *hashReader) Sum8() uint8 {
	return hr.Hash.(Hash8).Sum8()
}

// Sum16 returns the 16-bit checksum of the hash. It panics if hr.Hash doesn't
// implement the interface Hash16.
func (hr *hashReader) Sum16() uint16 {
	return hr.Hash.(Hash16).Sum16()
}