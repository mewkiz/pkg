// Package errorsutil implements some errors utility functions.
package errorsutil

import (
	"errors"
	"fmt"
	"path"
	"runtime"

	"github.com/mewkiz/pkg/term"
)

// New returns a new error string using the following format:
//
//	pkg.func (file:line): error: text
func New(text string) error {
	return backendNew(text, 2)
}

// backendNew returns a new error string using the following format:
//
//	pkg.func (file:line): error: text
func backendNew(text string, skip int) error {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return errors.New(text)
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return fmt.Errorf("(%s:%d): %s", path.Base(file), line, text)
	}
	return fmt.Errorf("%s (%s:%d): %s", f.Name(), path.Base(file), line, text)
}

// NewColor returns a new colorful error string using the following format:
//
//	pkg.func (file:line): error: text
func NewColor(text string) error {
	return backendNewColor(text, 2)
}

// backendNewColor returns a new colorful error string using the following
// format:
//
//	pkg.func (file:line): error: text
func backendNewColor(text string, skip int) error {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return errors.New(term.RedBold("error: ") + text)
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		format := term.WhiteBold("(%s:%d): ") + term.RedBold("error: ") + "%s"
		return fmt.Errorf(format, path.Base(file), line, text)
	}
	format := term.MagentaBold("%s") + term.WhiteBold(" (%s:%d): ") + term.RedBold("error: ") + "%s"
	return fmt.Errorf(format, f.Name(), path.Base(file), line, text)
}

// Errorf returns a new error string, based on the provided format string, using
// the following format:
//
//	pkg.func (file:line): error: text
func Errorf(format string, a ...interface{}) error {
	text := fmt.Sprintf(format, a...)
	return backendNew(text, 2)
}

// ErrorfColor returns a new colorful error string, based on the provided format
// string, using the following format:
//
//	pkg.func (file:line): error: text
func ErrorfColor(format string, a ...interface{}) error {
	text := fmt.Sprintf(format, a...)
	return backendNewColor(text, 2)
}
