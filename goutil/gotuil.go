// Package goutil implements some golang relevant utility functions.
package goutil

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"
)

// SrcDir tries to locate dir in GOPATH/src/ or GOROOT/src/pkg/ and returns its
// full path and true if successful. GOPATH may contain a list of multiple
// paths.
func SrcDir(dir string) (absDir string, err error) {
	for _, srcDir := range build.Default.SrcDirs() {
		absDir = filepath.Join(srcDir, dir)
		_, err := os.Stat(absDir)
		if err == nil {
			return absDir, nil
		}
	}
	return "", fmt.Errorf("goutil.SrcDir: unable to locate directory (%q) in GOPATH/src/ (%q) or GOROOT/src/pkg/ (%q).", dir, os.Getenv("GOPATH"), os.Getenv("GOROOT"))
}

// AbsImpPath tries to locate the absolute import path based on the provided
// import path.
func AbsImpPath(relImpPath string) (absImpPath string, err error) {
	if !build.IsLocalImport(relImpPath) {
		return relImpPath, nil
	}
	absPath, err := filepath.Abs(relImpPath)
	if err != nil {
		return "", fmt.Errorf("goutil.AbsImpPath: %v", err)
	}
	for _, srcDir := range build.Default.SrcDirs() {
		if strings.HasPrefix(absPath, srcDir) {
			absImpPath := absPath[len(srcDir):]
			if len(absImpPath) > 0 && absImpPath[0] == '/' {
				absImpPath = absImpPath[1:]
			}
			return absImpPath, nil
		}
	}
	return "", fmt.Errorf("goutil.AbsImpPath: unable to locate absolute import path for %q", relImpPath)
}
