// Package stringsutil implements some strings utility functions.
package stringsutil

import "strings"
import "unicode"

// IndexAfter returns the index directly after the first instance of sep in s,
// or -1 if sep is not present in s.
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

// SplitCamelCase splits the string s at each run of upper case runes and
// returns and array of slices of s.
func SplitCamelCase(s string) (words []string) {
	fieldStart := 0
	for i, r := range s {
		if i != 0 && unicode.IsUpper(r) {
			words = append(words, s[fieldStart:i])
			fieldStart = i
		}
	}
	if fieldStart != len(s) {
		words = append(words, s[fieldStart:])
	}
	return words
}
