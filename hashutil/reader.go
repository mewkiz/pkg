package hashutil

import (
	"hash"
	"io"
)

// HashReader wraps the basic hash.Hash and io.Reader interfaces. All data read
// from the io.Reader is also added to the running hash, unless Hash is set to
// nil.
type HashReader struct {
	hash.Hash
	io.Reader
}

// NewHashReader returns a HashReader based on the provided io.Reader. All data
// read is also added to the running hash, unless Hash is set to nil.
func NewHashReader(r io.Reader, h hash.Hash) (hr *HashReader) {
	hr = &HashReader{
		Hash:   h,
		Reader: r,
	}
	return hr
}

// Read reads up to len(buf) bytes into buf. All data read is also added to the
// running hash, unless Hash is set to nil. It returns the number of bytes read
// (0 <= n <= len(buf)) and any error encountered.
func (hr *HashReader) Read(buf []byte) (n int, err error) {
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
func (hr *HashReader) Sum8() uint8 {
	return hr.Hash.(Hash8).Sum8()
}

// Sum16 returns the 16-bit checksum of the hash. It panics if hr.Hash doesn't
// implement the interface Hash16.
func (hr *HashReader) Sum16() uint16 {
	return hr.Hash.(Hash16).Sum16()
}

// Sum32 returns the 32-bit checksum of the hash. It panics if hr.Hash doesn't
// implement the interface hash.Hash32.
func (hr *HashReader) Sum32() uint32 {
	return hr.Hash.(hash.Hash32).Sum32()
}

// Sum64 returns the 64-bit checksum of the hash. It panics if hr.Hash doesn't
// implement the interface hash.Hash64.
func (hr *HashReader) Sum64() uint64 {
	return hr.Hash.(hash.Hash64).Sum64()
}
