package errutil_test

import (
	"fmt"

	"github.com/mewkiz/pkg/errutil"
)

func ExampleNew() {
	errutil.UseColor = false
	err := errutil.New("failure.")
	fmt.Println(err)

	// Output:
	// github.com/mewkiz/pkg/errutil_test.ExampleNew (err_test.go:11): failure.
}
