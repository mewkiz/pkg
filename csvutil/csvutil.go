// Package csvutil implements CSV utility functions.
package csvutil

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"

	"github.com/jszwec/csvutil"
	"github.com/pkg/errors"
)

// Parse parses the given CSV stream into v.
func Parse(r io.Reader, v interface{}) error {
	br := bufio.NewReader(r)
	rr := csv.NewReader(br)
	dec, err := csvutil.NewDecoder(rr)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := dec.Decode(v); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// ParseFile parses the given CSV file into v.
func ParseFile(csvPath string, v interface{}) (err error) {
	f, err := os.Open(csvPath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			err = errors.WithStack(e)
		}
	}()
	if err := Parse(f, v); err != nil {
		return errors.Wrapf(err, "unable to parse %q", csvPath)
	}
	return err
}

// Write marshals v into CSV format, writing to w.
func Write(w io.Writer, v interface{}) error {
	ww := csv.NewWriter(w)
	enc := csvutil.NewEncoder(ww)
	if err := enc.Encode(v); err != nil {
		return errors.WithStack(err)
	}
	ww.Flush() // flush pending writes.
	if err := ww.Error(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// WriteFile marshals v into CSV format, writing to csvPath.
func WriteFile(csvPath string, v interface{}) (err error) {
	fd, err := os.Create(csvPath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		if e := fd.Close(); e != nil {
			err = errors.WithStack(e)
		}
	}()
	if err := Write(fd, v); err != nil {
		return errors.WithStack(err)
	}
	return err
}
