// Package goutil implements some golang relevant utility functions.
package goutil

import "fmt"
import "os"
import "path/filepath"

// SrcDir tries to locate dir in GOPATH/src/ and returns it's full path and true
// if successful. GOPATH may contain a list of multiple paths.
func SrcDir(dir string) (srcDir string, err error) {
	for _, goPath := range filepath.SplitList(os.Getenv("GOPATH")) {
		srcDir = fmt.Sprintf("%s/src/%s", goPath, dir)
		_, err := os.Stat(srcDir)
		if err == nil {
			return srcDir, nil
		}
	}
	return "", fmt.Errorf("goutil.SrcDir: unable to locate directory (%q) in GOPATH/src/ (%q).", dir, os.Getenv("GOPATH"))
}
