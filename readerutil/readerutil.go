// Package readerutil implements io.Reader utility functions.
package readerutil

import "io"

// ReadByte reads and returns one byte from the provided io.Reader.
func ReadByte(r io.Reader) (b byte, err error) {
	buf := make([]byte, 1)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return 0, err
	}

	return buf[0], nil
}
