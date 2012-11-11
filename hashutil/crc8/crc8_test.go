package crc8_test

import "io"
import "testing"

import "github.com/mewkiz/pkg/hashutil/crc8"

type test struct {
	want uint8
	in   string
}

var golden = []test{
	{0x00, ""},
	{0x20, "a"},
	{0xC9, "ab"},
	{0x5F, "abc"},
	{0xC1, "The quick brown fox jumps over the lazy dog"},
}

func TestGolden(t *testing.T) {
	for _, g := range golden {
		h := crc8.NewATM()
		io.WriteString(h, g.in)
		got := h.Sum8()
		if got != g.want {
			t.Errorf("ATM(%q); expected 0x%02X, got 0x%02X.", g.in, g.want, got)
		}
	}
}

func BenchmarkCrc8KB(b *testing.B) {
	b.SetBytes(1024)
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i)
	}
	h := crc8.NewATM()
	in := make([]byte, 0, h.Size())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Reset()
		h.Write(data)
		h.Sum(in)
	}
}
