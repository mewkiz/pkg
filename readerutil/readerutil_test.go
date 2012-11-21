package readerutil

import "encoding/binary"
import "os"
import "testing"

type testSize struct {
	path string
	want int64 // size in bytes.
}

func TestSize(t *testing.T) {
	golden := []testSize{
		{path: "testdata/utf16be.txt", want: 54},
		{path: "testdata/utf16le.txt", want: 54},
		{path: "testdata/utf8.txt", want: 31},
	}
	for _, g := range golden {
		fr, err := os.Open(g.path)
		if err != nil {
			t.Fatal(err)
		}
		defer fr.Close()

		// Verify file size.
		got, err := Size(fr)
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
		{path: "testdata/utf16be.txt", want: false},
		{path: "testdata/utf16le.txt", want: false},
		{path: "testdata/utf8.txt", want: true},
	}
	for _, g := range golden {
		fr, err := os.Open(g.path)
		if err != nil {
			t.Fatal(err)
		}
		defer fr.Close()

		// Verify file encoding.
		got, err := IsUTF8(fr)
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
		{path: "testdata/utf16be.txt", order: binary.BigEndian, want: true},
		{path: "testdata/utf16be.txt", order: binary.LittleEndian, want: false},
		{path: "testdata/utf16le.txt", order: binary.BigEndian, want: false},
		{path: "testdata/utf16le.txt", order: binary.LittleEndian, want: true},
		{path: "testdata/utf8.txt", order: binary.BigEndian, want: false},
		{path: "testdata/utf8.txt", order: binary.LittleEndian, want: false},
	}
	for _, g := range golden {
		fr, err := os.Open(g.path)
		if err != nil {
			t.Fatal(err)
		}
		defer fr.Close()

		// Verify file encoding.
		got, err := IsUTF16(fr, g.order)
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
