// Package osutil implements some os utility functions.
package osutil

import (
	"fmt"
	"os"
)

// Exists reports whether the given file or directory exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	panic(fmt.Errorf("unable to stat path %q; %v", path, err))
}
