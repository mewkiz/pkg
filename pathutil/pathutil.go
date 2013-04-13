// Package pathutil implements path utility functions.
package pathutil

import (
	"path"
)

// TrimExt returns a slice of the string filePath without the extension.
func TrimExt(filePath string) string {
	ext := path.Ext(filePath)
	return filePath[:len(filePath)-len(ext)
}
