// Package goutil implements some golang relevant utility functions.
package goutil

import (
	"go/build"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// SrcDir tries to locate dir in GOPATH/src/ or GOROOT/src/pkg/ and returns its
// full path. GOPATH may contain a list of paths.
func SrcDir(dir string) (absDir string, err error) {
	for _, srcDir := range build.Default.SrcDirs() {
		absDir = filepath.Join(srcDir, dir)
		finfo, err := os.Stat(absDir)
		if err == nil && finfo.IsDir() {
			return absDir, nil
		}
	}
	return "", errors.Errorf("unable to locate directory (%q) in GOPATH/src/ (%q) or GOROOT/src/pkg/ (%q)", dir, os.Getenv("GOPATH"), os.Getenv("GOROOT"))
}

// AbsImpPath tries to locate the absolute import path based on the provided
// import path.
func AbsImpPath(relImpPath string) (absImpPath string, err error) {
	if !build.IsLocalImport(relImpPath) {
		return relImpPath, nil
	}
	absPath, err := filepath.Abs(relImpPath)
	if err != nil {
		return "", errors.WithStack(err)
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
	return "", errors.Errorf("unable to locate absolute import path for %q", relImpPath)
}

// RelImpPath tries to locate the relative import path of the provided path. A
// pre-condition is that the given path resides within GOPATH.
func RelImpPath(path string) (relImpPath string, err error) {
	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)
		if err != nil {
			return "", errors.WithStack(err)
		}
	}
	for _, srcDir := range build.Default.SrcDirs() {
		if filepath.HasPrefix(path, srcDir) {
			return path[len(srcDir)+1:], nil
		}
	}
	return "", errors.Errorf("unable to locate path (%q) in GOPATH/src/ (%q) or GOROOT/src/pkg/ (%q)", path, os.Getenv("GOPATH"), os.Getenv("GOROOT"))
}
