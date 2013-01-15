package crc16_test

import "io"
import "testing"

import crc16 "."

type test struct {
	want uint16
	in   string
}

var golden = []test{
	{0x0000, ""},
	{0x8145, "a"},
	{0xC749, "ab"},
	{0xCADB, "abc"},
	{0x60AE, "The quick brown fox jumps over the lazy dog"},
}

func TestGolden(t *testing.T) {
	for _, g := range golden {
		h := crc16.NewIBM()
		io.WriteString(h, g.in)
		got := h.Sum16()
		if got != g.want {
			t.Errorf("IBM(%q); expected 0x%02X, got 0x%02X.", g.in, g.want, got)
		}
	}
}

func BenchmarkCrc16KB(b *testing.B) {
	b.SetBytes(1024)
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i)
	}
	h := crc16.NewIBM()
	in := make([]byte, 0, h.Size())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Reset()
		h.Write(data)
		h.Sum(in)
	}
}
