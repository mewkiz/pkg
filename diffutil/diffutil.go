// Package diffutil displays differences using Git.
package diffutil

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

// Diff displays the difference between a and b using Git.
func Diff(a, b string, words bool, filename string) error {
	dir, err := ioutil.TempDir("/tmp", "diff_")
	if err != nil {
		return errors.WithStack(err)
	}
	filename = filepath.Base(filename)
	if len(filename) == 0 {
		filename = "foo"
	}
	path := filepath.Join(dir, filename)
	if err := ioutil.WriteFile(path, []byte(a), 0644); err != nil {
		return errors.WithStack(err)
	}
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return errors.WithStack(err)
	}
	cmd = exec.Command("git", "add", path)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return errors.WithStack(err)
	}
	if err := ioutil.WriteFile(path, []byte(b), 0644); err != nil {
		return errors.WithStack(err)
	}
	if words {
		cmd = exec.Command("git", "diff", "--color-words")
	} else {
		cmd = exec.Command("git", "diff")
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return errors.WithStack(err)
	}
	if err := os.RemoveAll(dir); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
