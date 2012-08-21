// Package stringsutil implements some strings utility functions.
package stringsutil

import "strings"

// Index returns the index directly after the first instance of sep in s, or -1
// if sep is not present in s.
func IndexAfter(s, sep string) int {
	pos := strings.Index(s, sep)
	if pos == -1 {
		return -1
	}
	return pos + len(sep)
}

// Reverse returns a reversed version of s.
func Reverse(s string) (rev string) {
	for _, r := range s {
		rev = string(r) + rev
	}
	return rev
}
