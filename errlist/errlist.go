// Package errlist provides error list handling primitives.
package errlist

import "strings"

// Errors represents a list of errors, and implements the error interface.
type Errors []error

// Error returns a string representation of the list of errors.
func (es Errors) Error() string {
	if len(es) < 1 {
		return ""
	}
	var ss []string
	for _, e := range es {
		ss = append(ss, e.Error())
	}
	return strings.Join(ss, "; ")
}
