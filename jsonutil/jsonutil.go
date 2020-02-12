// Package jsonutil implements JSON utility functions.
package jsonutil

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/pkg/errors"
)

// Parse parses the given JSON stream into v.
func Parse(r io.Reader, v interface{}) error {
	br := bufio.NewReader(r)
	dec := json.NewDecoder(br)
	if err := dec.Decode(v); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// ParseFile parses the given JSON file into v.
func ParseFile(jsonPath string, v interface{}) error {
	f, err := os.Open(jsonPath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	if err := Parse(f, v); err != nil {
		return errors.Wrapf(err, "unable to parse %q", jsonPath)
	}
	return nil
}

// Write marshals v into JSON format, writing to w.
func Write(w io.Writer, v interface{}) error {
	buf, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return errors.WithStack(err)
	}
	if _, err := w.Write(buf); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// WriteFile marshals v into JSON format, writing to jsonPath.
func WriteFile(jsonPath string, v interface{}) error {
	fd, err := os.Create(jsonPath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer fd.Close()
	return Write(fd, v)
}
