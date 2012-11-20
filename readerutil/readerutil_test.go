package readerutil

import "os"
import "testing"

type testSize struct {
	path string
	want int64 // size in bytes.
}

func TestSize(t *testing.T) {
	golden := []testSize{
		{path: "testdata/utf8.txt", want: 31},
		{path: "testdata/utf16be.txt", want: 54},
		{path: "testdata/utf16le.txt", want: 54},
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
