package bufioutil

import (
	"bufio"
	"io"
)

// Writer implements buffering for an io.Writer object.
type Writer struct {
	*bufio.Writer
}

// NewWriter returns a new Writer.
func NewWriter(r io.Writer) (br Writer) {
	br.Writer = bufio.NewWriter(r)
	return br
}

// WriteLine writes a single line, including an newline character.
func (bw Writer) WriteLine(str string) (i int, err error) {
	i, err = bw.Writer.WriteString(str + "\n")
	if err != nil && i < len(str) {
		return i, err
	}
	bw.Writer.Flush()
	return i, nil
}
