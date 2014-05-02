package readerutil

import (
	"io"
	"strings"
)

// A filter discards characters from charset when reading from the underlying
// io.Reader.
type filter struct {
	io.Reader
	charset string
}

// NewFilter returns an io.Reader which discards characters from charset when
// reading from r.
func NewFilter(r io.Reader, charset string) io.Reader {
	f := &filter{
		Reader:  r,
		charset: charset,
	}
	return f
}

const (
	// Whitespace characters include space, tab, carriage return and newline.
	Whitespace = " \t\r\n"
)

// NewSpaceFilter returns an io.Reader which discards whitespace characters when
// reading from r.
func NewSpaceFilter(r io.Reader) io.Reader {
	f := &filter{
		Reader:  r,
		charset: Whitespace,
	}
	return f
}

func (f *filter) Read(buf []byte) (n int, err error) {
	m, err := f.Reader.Read(buf)
	for i := 0; i < m; i++ {
		// Ignore whitespace.
		if strings.IndexByte(f.charset, buf[i]) != -1 {
			continue
		}

		if n != i {
			buf[n] = buf[i]
		}
		n++
	}
	return n, err
}
