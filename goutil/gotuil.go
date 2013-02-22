// Package goutil implements some golang relevant utility functions.
package goutil

import "os"
import "path/filepath"

// SrcDir tries to locate dir in GOPATH/src/ and returns it's full path and true
// if successful. GOPATH may contain a list of multiple paths.
func SrcDir(dir string) (srcDir string, ok bool) {
	for _, goPath := range filepath.SplitList(os.Getenv("GOPATH")) {
		srcDir = goPath + "src/" + dir
		_, err := os.Stat(srcDir)
		if err == nil {
			return srcDir, true
		}
	}
	return "", false
}
