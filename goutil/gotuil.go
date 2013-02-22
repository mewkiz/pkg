// Package goutil implements some golang relevant utility functions.
package goutil

import "os"
import "path/filepath"

// Work around to find the source directory when GOPATH contains a list of multiple paths.
func SrcDir(repo string) (srcDir string, ok bool) {
	for _, goPath := range filepath.SplitList(os.Getenv("GOPATH")) {
		baseDir = goPath + srcDir
		_, err := os.Stat(baseDir)
		if err == nil {
			return baseDir, ok
		}
	}
	return "", false
}
