package readerutil_test

import "encoding/binary"
import "io"
import "os"
import "testing"

import readerutil "."

type testSize struct {
	path string
	want int64 // size in bytes.
}

func TestSize(t *testing.T) {
	golden := []testSize{
		{path: "testdata/utf16be_crlf.txt", want: 62},
		{path: "testdata/utf16be.txt", want: 56},
		{path: "testdata/utf16le_crlf.txt", want: 62},
		{path: "testdata/utf16le.txt", want: 56},
		{path: "testdata/utf8_crlf.txt", want: 35},
		{path: "testdata/utf8.txt", want: 32},
	}
	for _, g := range golden {
		fr, err := os.Open(g.path)
		if err != nil {
			t.Error(err)
			continue
		}
		defer fr.Close()

		// Verify file size.
		got, err := readerutil.Size(fr)
		if err != nil {
			t.Error(err)
			continue
		}
		if got != g.want {
			t.Errorf("%s: expected %d, got %d.", g.path, g.want, got)
			continue
		}
	}
}

type testIsUTF8 struct {
	path string
	want bool
}

func TestIsUTF8(t *testing.T) {
	golden := []testIsUTF8{
		{path: "testdata/utf16be_crlf.txt", want: false},
		{path: "testdata/utf16be.txt", want: false},
		{path: "testdata/utf16le_crlf.txt", want: false},
		{path: "testdata/utf16le.txt", want: false},
		{path: "testdata/utf8_crlf.txt", want: true},
		{path: "testdata/utf8.txt", want: true},
	}
	for _, g := range golden {
		fr, err := os.Open(g.path)
		if err != nil {
			t.Error(err)
			continue
		}
		defer fr.Close()

		// Verify file encoding.
		got, err := readerutil.IsUTF8(fr)
		if err != nil {
			t.Error(err)
			continue
		}
		if got != g.want {
			t.Errorf("%s: expected %t, got %t.", g.path, g.want, got)
			continue
		}
	}
}

type testIsUTF16 struct {
	path  string
	order binary.ByteOrder
	want  bool
}

func TestIsUTF16(t *testing.T) {
	golden := []testIsUTF16{
		{path: "testdata/utf16be_crlf.txt", order: binary.BigEndian, want: true},
		{path: "testdata/utf16be_crlf.txt", order: binary.LittleEndian, want: false},
		{path: "testdata/utf16be.txt", order: binary.BigEndian, want: true},
		{path: "testdata/utf16be.txt", order: binary.LittleEndian, want: false},
		{path: "testdata/utf16le_crlf.txt", order: binary.BigEndian, want: false},
		{path: "testdata/utf16le_crlf.txt", order: binary.LittleEndian, want: true},
		{path: "testdata/utf16le.txt", order: binary.BigEndian, want: false},
		{path: "testdata/utf16le.txt", order: binary.LittleEndian, want: true},
		{path: "testdata/utf8_crlf.txt", order: binary.BigEndian, want: false},
		{path: "testdata/utf8_crlf.txt", order: binary.LittleEndian, want: false},
		{path: "testdata/utf8.txt", order: binary.BigEndian, want: false},
		{path: "testdata/utf8.txt", order: binary.LittleEndian, want: false},
	}
	for _, g := range golden {
		fr, err := os.Open(g.path)
		if err != nil {
			t.Error(err)
			continue
		}
		defer fr.Close()

		// Verify file encoding.
		got, err := readerutil.IsUTF16(fr, g.order)
		if err != nil {
			t.Error(err)
			continue
		}
		if got != g.want {
			t.Errorf("%s: expected %t, got %t.", g.path, g.want, got)
			continue
		}
	}
}

type testNewLineReader struct {
	path  string
	lines []string
}

func TestNewLineReader(t *testing.T) {
	golden := []testNewLineReader{
		{
			path: "testdata/utf16be_crlf.txt",
			lines: []string{
				"testing 123",
				"",
				"hello 世界",
				"mewmew",
			},
		},
		{
			path: "testdata/utf16be.txt",
			lines: []string{
				"testing 123",
				"",
				"hello 世界",
				"mewmew",
			},
		},
		{
			path: "testdata/utf16le_crlf.txt",
			lines: []string{
				"testing 123",
				"",
				"hello 世界",
				"mewmew",
			},
		},
		{
			path: "testdata/utf16le.txt",
			lines: []string{
				"testing 123",
				"",
				"hello 世界",
				"mewmew",
			},
		},
		{
			path: "testdata/utf8_crlf.txt",
			lines: []string{
				"testing 123",
				"",
				"hello 世界",
				"mewmew",
			},
		},
		{
			path: "testdata/utf8.txt",
			lines: []string{
				"testing 123",
				"",
				"hello 世界",
				"mewmew",
			},
		},
	}
	for _, g := range golden {
		fr, err := os.Open(g.path)
		if err != nil {
			t.Error(err)
			continue
		}
		defer fr.Close()

		// Verify lines.
		lr, err := readerutil.NewLineReader(fr)
		if err != nil {
			t.Error(err)
			continue
		}
		lineNum := 0
		for {
			got, err := lr.ReadLine()
			if err != nil {
				if err != io.EOF {
					t.Error(err)
				}
				// break on io.EOF
				break
			}
			if lineNum >= len(g.lines) {
				t.Errorf("lines slice out of bounds (%d >= %d).", lineNum, len(g.lines))
				break
			}
			want := g.lines[lineNum]
			lineNum++
			if got != want {
				t.Errorf("%s: expected %q, got %q.", g.path, want, got)
				continue
			}
		}
		if lineNum != len(g.lines) {
			t.Errorf("tested %d out of %d lines.", lineNum, len(g.lines))
		}
	}
}
