// Package bitutil implements common binary encoding and decoding algorithms.
package bitutil

import (
	"github.com/eaburns/bit"
)

// DecodeUnary decodes an unary coded integer and returns it.
//
// Examples of unary coded binary on the left and decoded decimal on the right:
//
//    1       => 1
//    01      => 2
//    001     => 3
//    0001    => 4
//    00001   => 5
//    000001  => 6
//    0000001 => 7
func DecodeUnary(br *bit.Reader) (n uint, err error) {
	for {
		bit, err := br.Read(1)
		if err != nil {
			return 0, err
		}
		n++
		if bit == 1 {
			break
		}
	}
	return n, nil
}
