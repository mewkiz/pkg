package bufioutil

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"
)

func TestReadLineBufferBoundaries(t *testing.T) {
	for _, length := range []int{0, 1, 4095, 4096, 4097, 8192} {
		for _, ending := range []string{"", "\n", "\r\n"} {
			t.Run(fmt.Sprintf("%d/%q", length, ending), func(t *testing.T) {
				want := strings.Repeat("x", length)
				reader := NewReader(strings.NewReader(want + ending))
				var wantErr error
				if length == 0 && ending == "" {
					wantErr = io.EOF
				}
				got, err := reader.ReadLine()
				if got != want || err != wantErr {
					t.Fatalf("ReadLine() returned %d bytes and %v; want %d bytes and %v", len(got), err, len(want), wantErr)
				}
				if got, err := reader.ReadLine(); got != "" || err != io.EOF {
					t.Fatalf("next ReadLine() = (%q, %v), want (empty, EOF)", got, err)
				}
			})
		}
	}
}

func TestReadLineErrorAfterPrefix(t *testing.T) {
	reader := NewReader(io.MultiReader(
		strings.NewReader(strings.Repeat("x", 4096)),
		iotest.ErrReader(io.ErrUnexpectedEOF),
	))
	if got, err := reader.ReadLine(); got != "" || err != io.ErrUnexpectedEOF {
		t.Fatalf("ReadLine() = (%q, %v), want (empty, unexpected EOF)", got, err)
	}
}

func TestReadLinesFinalFullBuffer(t *testing.T) {
	want := []string{"first", strings.Repeat("x", 4096)}
	got, err := ReadLines(strings.NewReader(strings.Join(want, "\n")))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ReadLines() returned %d lines, want %d", len(got), len(want))
	}
}
