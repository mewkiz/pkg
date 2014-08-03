// Package bit provides bit reading operations for io.Readers.
package bit

import (
	"fmt"
	"io"
)

// A Reader handles bit reading operations. It buffers bits up to the next byte
// boundary.
type Reader struct {
	// Underlying reader.
	r io.Reader
	// Temporary read buffer.
	buf [8]uint8
	// Between 0 and 7 buffered bits since previous read operations.
	x uint8
	// The number of buffered bits in x.
	n uint
}

// NewReader returns a new Reader that reads bits from r.
func NewReader(r io.Reader) (br *Reader) {
	return &Reader{r: r}
}

// Read reads and returns the next n bits, at most 64. It buffers bits up to the
// next byte boundary.
func (br *Reader) Read(n uint) (x uint64, err error) {
	if n > 64 {
		return 0, fmt.Errorf("bit.Reader.Read: invalid number of bits; n (%d) exceeds 64", n)
	}

	// Read buffered bits.
	if br.n > 0 {
		switch {
		case br.n == n:
			br.n = 0
			return uint64(br.x), nil
		case br.n > n:
			br.n -= n
			mask := ^uint8(0) << br.n
			x = uint64(br.x&mask) >> br.n
			br.x &^= mask
			return x, nil
		}
		n -= br.n
		x = uint64(br.x)
		br.n = 0
	}

	// Fill the temporary buffer.
	bytes := n / 8
	bits := n % 8
	if bits > 0 {
		bytes++
	}
	_, err = io.ReadFull(br.r, br.buf[:bytes])
	if err != nil {
		return 0, err
	}

	// Read bits from the temporary buffer.
	for _, b := range br.buf[:bytes-1] {
		x <<= 8
		x |= uint64(b)
	}
	b := br.buf[bytes-1]
	if bits > 0 {
		x <<= bits
		br.n = 8 - bits
		mask := ^uint8(0) << br.n
		x |= uint64(b&mask) >> br.n
		br.x = b & ^mask
	} else {
		x <<= 8
		x |= uint64(b)
	}

	return x, nil
}

// ReadFields reads a series of consecutive bit sequences as specified by the
// argument list.
func (br *Reader) ReadFields(ns ...uint) (xs []uint64, err error) {
	xs = make([]uint64, 0, len(ns))
	for i, n := range ns {
		x, err := br.Read(n)
		if err != nil {
			if err == io.EOF && i != 0 {
				return xs, io.ErrUnexpectedEOF
			}
			return xs, err
		}
		xs = append(xs, x)
	}
	return xs, nil
}

// Buffered returns the number of bits that can be read from the current buffer,
// in range [0, 7], which also represent the remaining number of bits to the
// next byte boundary.
func (br *Reader) Buffered() uint {
	return br.n
}
