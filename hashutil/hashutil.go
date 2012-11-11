// Package hashutil provides interfaces for hash functions.
package hashutil

import "hash"

// Hash8 is the common interface implemented by all 8-bit hash functions.
type Hash8 interface {
	hash.Hash
	Sum8() uint8
}
