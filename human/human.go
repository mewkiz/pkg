// Package human implements human representations of various data.
package human

import (
	"strconv"
)

// Ordinal returns a human representation of n in the form of 1st, 2nd, 3rd,
// 4th, etc...
func Ordinal(n int) string {
	var suffix string
	switch n % 10 {
	case 1:
		suffix = "st"
	case 2:
		suffix = "nd"
	case 3:
		suffix = "rd"
	default:
		suffix = "th"
	}
	return strconv.Itoa(n) + suffix
}
