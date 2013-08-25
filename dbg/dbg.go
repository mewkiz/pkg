// Package dbg implements formatted I/O which can be enabled or disabled at
// runtime.
package dbg

import (
	"fmt"
)

// When Debug is true, output is enabled.
var Debug = true

// Print formats using the default formats for its operands and writes to
// standard output. Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func Print(a ...interface{}) (n int, err error) {
	if !Debug {
		return 0, nil
	}
	return fmt.Print(a...)
}

// Printf formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
func Printf(format string, a ...interface{}) (n int, err error) {
	if !Debug {
		return 0, nil
	}
	return fmt.Printf(format, a...)
}

// Println formats using the default formats for its operands and writes to
// standard output. Spaces are always added between operands and a newline is
// appended. It returns the number of bytes written and any write error
// encountered.
func Println(a ...interface{}) (n int, err error) {
	if !Debug {
		return 0, nil
	}
	return fmt.Println(a...)
}
