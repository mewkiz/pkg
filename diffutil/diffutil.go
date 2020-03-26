// Package diffutil displays differences using Git.
package diffutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/d4l3k/messagediff"
	"github.com/mewkiz/pkg/natsort"
	"github.com/mewkiz/pkg/term"
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
		cmd = exec.Command("git", "diff", "--color-words", "--color=always")
	} else {
		cmd = exec.Command("git", "diff", "--color=always")
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

// PrettyDiff returns a pretty-printed colour output of the deep diff of a and
// b. The boolean return value indicates whether a and b are equal.
func PrettyDiff(a, b interface{}) (string, bool) {
	diff, equal := messagediff.DeepDiff(a, b)
	type T struct {
		path string
		s    string
	}
	var ts []T
	for path, added := range diff.Added {
		s := term.Green(fmt.Sprintf("+ %s = %#v\n", path, added))
		ts = append(ts, T{path: path.String(), s: s})
	}
	for path, removed := range diff.Removed {
		s := term.Red(fmt.Sprintf("- %s = %#v\n", path, removed))
		ts = append(ts, T{path: path.String(), s: s})
	}
	for path, modified := range diff.Modified {
		old := term.Red(fmt.Sprintf("- %s = %#v\n", path, modified.Old))
		new := term.Green(fmt.Sprintf("+ %s = %#v\n", path, modified.New))
		s := old + new
		ts = append(ts, T{path: path.String(), s: s})
	}
	sort.Slice(ts, func(i, j int) bool {
		return natsort.Less(ts[i].path, ts[j].path)
	})
	var ss []string
	for _, t := range ts {
		ss = append(ss, t.s)
	}
	return strings.Join(ss, ""), equal
}
