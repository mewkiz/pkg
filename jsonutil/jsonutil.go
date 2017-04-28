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
	return dec.Decode(v)
}

// ParseFile parses the given JSON file into v.
func ParseFile(path string, v interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	br := bufio.NewReader(f)
	dec := json.NewDecoder(br)
	return dec.Decode(v)
}
