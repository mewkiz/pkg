// Package bufioutil implements utility functions for buffered I/O.
package bufioutil

import (
	"bufio"
	"io"
	"os"
)

// Reader implements buffering for an io.Reader object.
type Reader struct {
	*bufio.Reader
}

// NewReader returns a new Reader.
func NewReader(r io.Reader) (br Reader) {
	br = Reader{
		Reader: bufio.NewReader(r),
	}
	return br
}

// ReadLine reads returns a single line from br, not including the end-of-line
// bytes.
func (br Reader) ReadLine() (line string, err error) {
	for {
		buf, isPrefix, err := br.Reader.ReadLine()
		if err != nil {
			return "", err
		}
		line += string(buf)
		if !isPrefix {
			break
		}
	}
	return line, nil
}

// ReadLines reads and returns all lines from br, not including the end-of-line
// bytes.
func (br Reader) ReadLines() (lines []string, err error) {
	for {
		line, err := br.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		lines = append(lines, line)
	}
	return lines, nil
}

// ReadLines reads from r and returns all lines, not including the end-of-line
// bytes.
func ReadLines(r io.Reader) (lines []string, err error) {
	br := NewReader(r)
	return br.ReadLines()
}

// LoadLines loads the provided file and returns all lines, not including the
// end-of-line bytes.
func LoadLines(filePath string) (lines []string, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	br := NewReader(f)
	return br.ReadLines()
}
