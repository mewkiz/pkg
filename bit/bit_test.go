// © 2013 the Bits Authors under the MIT license. See AUTHORS for the list of authors.
//
// Some benchmark functions in this file were adapted from github.com/bamiaux/iobit
// which came with the following copyright notice:
// Copyright 2013 Benoît Amiaux. All rights reserved.

package bit

import (
	"bytes"
	"io"
	"math/rand"
	"testing"
)

func TestRead(t *testing.T) {
	tests := []struct {
		data []byte
		ns   []uint
		vals []uint64
	}{
		{[]byte{0xFF}, []uint{1, 1, 1, 1, 1, 1, 1, 1}, []uint64{1, 1, 1, 1, 1, 1, 1, 1}},
		{[]byte{0xFF}, []uint{2, 2, 2, 2}, []uint64{0x3, 0x3, 0x3, 0x3}},
		{[]byte{0xFF}, []uint{3, 3, 2}, []uint64{0x7, 0x7, 0x3}},
		{[]byte{0xFF}, []uint{4, 4}, []uint64{0xF, 0xF}},
		{[]byte{0xFF}, []uint{5, 3}, []uint64{0x1F, 0x7}},
		{[]byte{0xFF}, []uint{6, 2}, []uint64{0x3F, 0x3}},
		{[]byte{0xFF}, []uint{7, 1}, []uint64{0x7F, 0x1}},
		{[]byte{0xFF}, []uint{8}, []uint64{0xFF}},

		{[]byte{0xAA}, []uint{1, 1, 1, 1, 1, 1, 1, 1}, []uint64{1, 0, 1, 0, 1, 0, 1, 0}},
		{[]byte{0xAA}, []uint{2, 2, 2, 2}, []uint64{0x2, 0x2, 0x2, 0x2}},
		{[]byte{0xAA}, []uint{3, 3, 2}, []uint64{0x5, 0x2, 0x2}},
		{[]byte{0xAA}, []uint{4, 4}, []uint64{0xA, 0xA}},
		{[]byte{0xAA}, []uint{5, 3}, []uint64{0x15, 0x2}},
		{[]byte{0xAA}, []uint{6, 2}, []uint64{0x2A, 0x2}},
		{[]byte{0xAA}, []uint{7, 1}, []uint64{0x55, 0x0}},
		{[]byte{0xAA}, []uint{8}, []uint64{0xAA}},

		{
			[]byte{0xAA, 0x55},
			[]uint{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			[]uint64{1, 0, 1, 0, 1, 0, 1, 0, 0, 1, 0, 1, 0, 1, 0, 1},
		},

		{
			[]byte{0xAA, 0x55},
			[]uint{7, 8, 1},
			[]uint64{0x55, 0x2A, 0x1},
		},

		{
			[]byte{0xAA, 0x55},
			[]uint{3, 3, 3, 3, 3, 1},
			[]uint64{0x5, 0x2, 0x4, 0x5, 0x2, 0x1},
		},

		{
			[]byte{0xAA, 0x55},
			[]uint{16},
			[]uint64{0xAA55},
		},

		{
			[]byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55},
			[]uint{32, 32},
			[]uint64{0xAA55AA55, 0xAA55AA55},
		},

		{
			[]byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55},
			[]uint{33, 31},
			[]uint64{0x154AB54AB, 0x2A55AA55},
		},
	}

	for _, test := range tests {
		r := NewReader(bytes.NewReader(test.data))
		if len(test.ns) != len(test.vals) {
			panic("Number of reads does not match number of results")
		}
		for i, n := range test.ns {
			m, err := r.Read(n)
			if err != nil {
				panic("Unexpected error: " + err.Error())
			}
			if m != test.vals[i] {
				t.Errorf("%v with reads %v: read %d gave %x, expected %x", test.data, test.ns, i, m, test.vals[i])
			}
		}
	}
}

func TestReadEOF(t *testing.T) {
	tests := []struct {
		data []byte
		n    uint
		err  error
	}{
		{[]byte{0xFF}, 8, nil},
		{[]byte{0xFF}, 2, nil},
		{[]byte{0xFF}, 9, io.ErrUnexpectedEOF},
		{[]byte{}, 1, io.EOF},
		{[]byte{0xFF, 0xFF}, 16, nil},
		{[]byte{0xFF, 0xFF}, 17, io.ErrUnexpectedEOF},
	}

	for _, test := range tests {
		r := NewReader(bytes.NewReader(test.data))
		if _, err := r.Read(test.n); err != test.err {
			t.Errorf("Reading %d from %v, expected err=%s, got err=%s", test.n, test.data, test.err, err)
		}
	}

}

func TestReadFields(t *testing.T) {
	tests := []struct {
		data []byte
		ns   []uint
		fs   []uint64
	}{
		{[]byte{0xFF}, []uint{1, 1, 1, 1, 1, 1, 1, 1}, []uint64{1, 1, 1, 1, 1, 1, 1, 1}},
		{[]byte{0xFF}, []uint{2, 2, 2, 2}, []uint64{0x3, 0x3, 0x3, 0x3}},
		{[]byte{0xFF}, []uint{3, 3, 2}, []uint64{0x7, 0x7, 0x3}},
		{[]byte{0xFF}, []uint{4, 4}, []uint64{0xF, 0xF}},
		{[]byte{0xFF}, []uint{5, 3}, []uint64{0x1F, 0x7}},
		{[]byte{0xFF}, []uint{6, 2}, []uint64{0x3F, 0x3}},
		{[]byte{0xFF}, []uint{7, 1}, []uint64{0x7F, 0x1}},
		{[]byte{0xFF}, []uint{8}, []uint64{0xFF}},

		{[]byte{0xAA}, []uint{1, 1, 1, 1, 1, 1, 1, 1}, []uint64{1, 0, 1, 0, 1, 0, 1, 0}},
		{[]byte{0xAA}, []uint{2, 2, 2, 2}, []uint64{0x2, 0x2, 0x2, 0x2}},
		{[]byte{0xAA}, []uint{3, 3, 2}, []uint64{0x5, 0x2, 0x2}},
		{[]byte{0xAA}, []uint{4, 4}, []uint64{0xA, 0xA}},
		{[]byte{0xAA}, []uint{5, 3}, []uint64{0x15, 0x2}},
		{[]byte{0xAA}, []uint{6, 2}, []uint64{0x2A, 0x2}},
		{[]byte{0xAA}, []uint{7, 1}, []uint64{0x55, 0x0}},
		{[]byte{0xAA}, []uint{8}, []uint64{0xAA}},

		{
			[]byte{0xAA, 0x55},
			[]uint{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			[]uint64{1, 0, 1, 0, 1, 0, 1, 0, 0, 1, 0, 1, 0, 1, 0, 1},
		},

		{
			[]byte{0xAA, 0x55},
			[]uint{7, 8, 1},
			[]uint64{0x55, 0x2A, 0x1},
		},

		{
			[]byte{0xAA, 0x55},
			[]uint{3, 3, 3, 3, 3, 1},
			[]uint64{0x5, 0x2, 0x4, 0x5, 0x2, 0x1},
		},

		{
			[]byte{0xAA, 0x55},
			[]uint{16},
			[]uint64{0xAA55},
		},

		{
			[]byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55},
			[]uint{32, 32},
			[]uint64{0xAA55AA55, 0xAA55AA55},
		},

		{
			[]byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55},
			[]uint{33, 31},
			[]uint64{0x154AB54AB, 0x2A55AA55},
		},
	}

	for _, test := range tests {
		r := NewReader(bytes.NewReader(test.data))
		if len(test.ns) != len(test.fs) {
			panic("Number of reads does not match number of results")
		}
		fs, err := r.ReadFields(test.ns...)
		if err != nil {
			panic("Unexpected error")
		}
		for i := range fs {
			if fs[i] != test.fs[i] {
				t.Errorf("Reading Fields %v from %v, expected %v, got %v", test.ns, test.data, test.ns, fs)
			}
		}
	}
}

func TestReadFieldsEOF(t *testing.T) {
	tests := []struct {
		data []byte
		ns   []uint
		err  error
	}{
		{[]byte{0xFF}, []uint{8}, nil},
		{[]byte{0xFF}, []uint{2}, nil},
		{[]byte{0xFF}, []uint{9}, io.ErrUnexpectedEOF},
		{[]byte{}, []uint{1}, io.EOF},
		{[]byte{0xFF, 0xFF}, []uint{16}, nil},
		{[]byte{0xFF, 0xFF}, []uint{17}, io.ErrUnexpectedEOF},
		{[]byte{0xFF}, []uint{1, 7}, nil},
		{[]byte{0xFF}, []uint{1, 8}, io.ErrUnexpectedEOF},
		{[]byte{}, []uint{1, 8}, io.EOF},
	}

	for _, test := range tests {
		r := NewReader(bytes.NewReader(test.data))
		if _, err := r.ReadFields(test.ns...); err != test.err {
			t.Errorf("Reading Fields %v from %v, expected err=%s, got err=%s", test.ns, test.data, test.err, err)
		}
	}

}

func BenchmarkReadAlign1(b *testing.B) {
	benchmarkReads(b, 64, 1)
}

func BenchmarkReadAlign32(b *testing.B) {
	benchmarkReads(b, 64, 32)
}

func BenchmarkReadAlign64(b *testing.B) {
	benchmarkReads(b, 64, 64)
}

func benchmarkReads(b *testing.B, chunk, align int) {
	size := 1 << 12
	buf, bits, _, last := prepareBenchmark(size, chunk, align)
	b.SetBytes(int64(len(buf)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := NewReader(bytes.NewReader(buf))
		for j := 0; j < last; j++ {
			r.Read(bits[j])
		}
	}
}

func prepareBenchmark(size, chunk, align int) ([]byte, []uint, []uint64, int) {
	buf := make([]byte, size)
	bits := make([]uint, size)
	values := make([]uint64, size)
	idx := 0
	last := 0
	for i := 0; i < size; i++ {
		val := getNumBits(idx, size*8, chunk, align)
		idx += val
		if val != 0 {
			last = i + 1
		}
		bits[i] = uint(val)
		values[i] = uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
	}
	return buf, bits, values, last
}

func getNumBits(read, max, chunk, align int) int {
	bits := 1
	if align != chunk {
		bits += rand.Intn(chunk / align)
	}
	bits *= align
	if read+bits > max {
		bits = max - read
	}
	if bits > chunk {
		panic("too many bits")
	}
	return bits
}
