package readerutil

import (
	"io"
	"os"
	"strings"
)

func ExampleNewFilter() {
	r := strings.NewReader("foo,bar;baz")
	f := NewFilter(r, ",;")
	io.Copy(os.Stdout, f)
	// Output: foobarbaz
}

func ExampleNewSpaceFilter() {
	r := strings.NewReader("foo bar\nbaz")
	f := NewSpaceFilter(r)
	io.Copy(os.Stdout, f)
	// Output: foobarbaz
}
