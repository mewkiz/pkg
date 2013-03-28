package bit

import "bytes"
import "errors"
import "io"
import "testing"

type testNewStream struct {
	buf  []byte
	want Stream
}

func TestNewStream(t *testing.T) {
	golden := []testNewStream{
		// empty original bitstream.
		{
			buf:  []byte{0x01, 0x02, 0x48, 0xFF, 0x19, 0x03, 0x05, 0x06, 0x07, 0x1D},
			want: Stream{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 1},
		},
		{
			buf:  []byte{0x01, 0x02, 0x48, 0xFF, 0x19, 0x03, 0x05, 0x06, 0x07},
			want: Stream{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1},
		},
		{
			buf:  []byte{0x01, 0x02, 0x48, 0xFF, 0x19, 0x03, 0x05},
			want: Stream{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 0, 1},
		},
		{
			buf:  []byte{0x01, 0x02, 0x48, 0xFF, 0x19},
			want: Stream{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 0, 0, 1},
		},
		{
			buf:  []byte{0x01, 0x02, 0x48, 0xFF},
			want: Stream{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1},
		},
		{
			buf:  []byte{0x01, 0x02, 0x48},
			want: Stream{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0},
		},
		{
			buf:  []byte{0x01},
			want: Stream{0, 0, 0, 0, 0, 0, 0, 1},
		},
		{
			buf:  []byte{},
			want: Stream{},
		},
	}

	for i, g := range golden {
		got := NewStream(g.buf)
		if !Equal(got, g.want) {
			t.Errorf("i %d: expected %q, got %q.", i, g.want, got)
		}
	}
}

type testStreamAppendBytes struct {
	buf  []byte
	src  Stream
	want Stream
}

func TestStreamAppendBytes(t *testing.T) {
	golden := []testStreamAppendBytes{
		// empty original bitstream.
		{
			buf:  []byte{0x01, 0x02, 0x48, 0xFF, 0x19, 0x03, 0x05, 0x06, 0x07, 0x1D},
			src:  Stream{},
			want: Stream{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 1},
		},
		{
			buf:  []byte{0x01, 0x02, 0x48, 0xFF, 0x19, 0x03, 0x05, 0x06, 0x07},
			src:  Stream{},
			want: Stream{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1},
		},
		{
			buf:  []byte{0x01, 0x02, 0x48, 0xFF, 0x19, 0x03, 0x05},
			src:  Stream{},
			want: Stream{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 0, 1},
		},
		{
			buf:  []byte{0x01, 0x02, 0x48, 0xFF, 0x19},
			src:  Stream{},
			want: Stream{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 0, 0, 1},
		},
		{
			buf:  []byte{0x01, 0x02, 0x48, 0xFF},
			src:  Stream{},
			want: Stream{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1},
		},
		{
			buf:  []byte{0x01, 0x02, 0x48},
			src:  Stream{},
			want: Stream{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0},
		},
		{
			buf:  []byte{0x01},
			src:  Stream{},
			want: Stream{0, 0, 0, 0, 0, 0, 0, 1},
		},
		{
			buf:  []byte{},
			src:  Stream{},
			want: Stream{},
		},

		// non-empty original bitstream.
		{
			buf:  []byte{0x01, 0x02, 0x48, 0xFF, 0x19, 0x03, 0x05, 0x06, 0x07, 0x1D},
			src:  Stream{0, 1, 1, 1, 1, 0, 1},
			want: Stream{0, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 1},
		},
		{
			buf:  []byte{0x01, 0x02, 0x48, 0xFF, 0x19, 0x03, 0x05, 0x06, 0x07},
			src:  Stream{0, 1, 1, 1, 1, 0, 1},
			want: Stream{0, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1},
		},
		{
			buf:  []byte{0x01, 0x02, 0x48, 0xFF, 0x19, 0x03, 0x05},
			src:  Stream{0, 1, 1, 1, 1, 0, 1},
			want: Stream{0, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 0, 1},
		},
		{
			buf:  []byte{0x01, 0x02, 0x48, 0xFF, 0x19},
			src:  Stream{0, 1, 1, 1, 1, 0, 1},
			want: Stream{0, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 0, 0, 1},
		},
		{
			buf:  []byte{0x01, 0x02, 0x48, 0xFF},
			src:  Stream{0, 1, 1, 1, 1, 0, 1},
			want: Stream{0, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1},
		},
		{
			buf:  []byte{0x01, 0x02, 0x48},
			src:  Stream{0, 1, 1, 1, 1, 0, 1},
			want: Stream{0, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0},
		},
		{
			buf:  []byte{0x01},
			src:  Stream{0, 1, 1, 1, 1, 0, 1},
			want: Stream{0, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1},
		},
		{
			buf:  []byte{},
			src:  Stream{0, 1, 1, 1, 1, 0, 1},
			want: Stream{0, 1, 1, 1, 1, 0, 1},
		},
	}

	for i, g := range golden {
		got := g.src.AppendBytes(g.buf)
		if !Equal(got, g.want) {
			t.Errorf("i %d: expected %q, got %q.", i, g.want, got)
		}
	}
}

type testStreamUint64 struct {
	bits Stream
	want uint64
}

func TestStreamUint64(t *testing.T) {
	golden := []testStreamUint64{
		{
			bits: Stream{0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1},
			want: 1073015,
		},
		{
			bits: Stream{0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1, 0, 1},
			want: 16765,
		},
		{
			bits: Stream{0, 1, 0, 0, 0, 0, 0, 1, 0, 1},
			want: 261,
		},
		{
			bits: Stream{0, 1, 0, 0, 0, 0, 0, 1},
			want: 65,
		},
		{
			bits: Stream{0, 1, 0, 0, 0, 0, 0},
			want: 32,
		},
		{
			bits: Stream{0, 1},
			want: 1,
		},
		{
			bits: Stream{0},
			want: 0,
		},
		{
			bits: Stream{},
			want: 0,
		},
	}

	for i, g := range golden {
		got := g.bits.Uint64()
		if got != g.want {
			t.Errorf("i %d: expected %v, got %v.", i, g.want, got)
		}
	}
}

type testStreamString struct {
	bits Stream
	want string
}

func TestStreamString(t *testing.T) {
	golden := []testStreamString{
		{
			bits: Stream{0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1},
			want: "01000001 01111101 110111",
		},
		{
			bits: Stream{0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1, 0, 1},
			want: "01000001 01111101",
		},
		{
			bits: Stream{0, 1, 0, 0, 0, 0, 0, 1, 0, 1},
			want: "01000001 01",
		},
		{
			bits: Stream{0, 1, 0, 0, 0, 0, 0, 1},
			want: "01000001",
		},
		{
			bits: Stream{0, 1, 0, 0, 0, 0, 0},
			want: "0100000",
		},
		{
			bits: Stream{0, 1},
			want: "01",
		},
		{
			bits: Stream{0},
			want: "0",
		},
		{
			bits: Stream{},
			want: "",
		},
	}

	for i, g := range golden {
		got := g.bits.String()
		if got != g.want {
			t.Errorf("i %d: expected %q, got %q.", i, g.want, got)
		}
	}
}

type testEqual struct {
	a, b Stream
	want bool
}

func TestEqual(t *testing.T) {
	golden := []testEqual{
		{
			// equal.
			a:    Stream{1, 0, 0, 1, 0, 0, 1, 0, 0, 1},
			b:    Stream{1, 0, 0, 1, 0, 0, 1, 0, 0, 1},
			want: true,
		},
		{
			// equal.
			a:    nil,
			b:    Stream{},
			want: true,
		},
		{
			// different length.
			a:    Stream{1, 0, 0, 1, 0, 0, 1, 0, 0},
			b:    Stream{1, 0, 0, 1, 0, 0, 1, 0, 0, 1},
			want: false,
		},
		{
			// different length.
			a:    Stream{1, 0, 0, 1, 0, 0, 1, 0, 0, 1},
			b:    Stream{1, 0, 0, 1, 0, 0, 1, 0, 0},
			want: false,
		},
		{
			// different length.
			a:    nil,
			b:    Stream{1, 0, 0, 1, 0, 0, 1, 0, 0},
			want: false,
		},
	}

	for i, g := range golden {
		got := Equal(g.a, g.b)
		if got != g.want {
			t.Errorf("i %d: expected %v, got %v; for:\na: %q\nb: %q.", i, g.want, got, g.a, g.b)
		}
	}
}

type testRead struct {
	r     io.ReadSeeker
	reads []testReadParams
}

type testReadParams struct {
	n    int
	bits Stream
	err  error
}

func TestRead(t *testing.T) {
	golden := []testRead{
		{
			// 10110111 01111011 11101111 11011111 11011111 11101111 11111011 11111111
			r: bytes.NewReader([]byte{0xB7, 0x7B, 0xEF, 0xDF, 0xDF, 0xEF, 0xFB, 0xFF}),
			reads: []testReadParams{
				{
					// before read: 10110111 01111011 11101111 11011111 11011111 11101111 11111011 11111111
					n:    1,
					bits: Stream{1},
					err:  nil,
				},
				{
					// before read: 0110111 01111011 11101111 11011111 11011111 11101111 11111011 11111111
					n:    6,
					bits: Stream{0, 1, 1, 0, 1, 1},
					err:  nil,
				},
				{
					// before read: 1 01111011 11101111 11011111 11011111 11101111 11111011 11111111
					n:    2,
					bits: Stream{1, 0},
					err:  nil,
				},
				{
					// before read: 1111011 11101111 11011111 11011111 11101111 11111011 11111111
					n:    11,
					bits: Stream{1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0},
					err:  nil,
				},
				{
					// before read: 1111 11011111 11011111 11101111 11111011 11111111
					n:    37,
					bits: Stream{1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1},
					err:  nil,
				},
				{
					// before read: 1111111
					n:    6,
					bits: Stream{1, 1, 1, 1, 1, 1},
					err:  nil,
				},
				{
					// before read: 1
					n:    1,
					bits: Stream{1},
					err:  nil,
				},
				{
					// before read: <empty>
					n:    1,
					bits: Stream{},
					err:  io.EOF,
				},
			},
		},

		{
			// 11101101 11010011 01111110 10110110
			r: bytes.NewReader([]byte{0xED, 0xD3, 0x7E, 0xB6}),
			reads: []testReadParams{
				{
					// before read: 11101101 11010011 01111110 10110110
					n:    33,
					bits: Stream{1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 0, 1, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 0, 1, 0, 1, 1, 0, 1, 1, 0},
					err:  io.ErrUnexpectedEOF,
				},
			},
		},

		{
			// 11101101 11010011 01111110 10110110
			r: bytes.NewReader([]byte{0xED, 0xD3, 0x7E, 0xB6}),
			reads: []testReadParams{
				{
					// before read: 11101101 11010011 01111110 10110110
					n:    17,
					bits: Stream{1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 0, 1, 0, 0, 1, 1, 0},
					err:  nil,
				},
				{
					// before read: 1111110 10110110
					n:    12,
					bits: Stream{1, 1, 1, 1, 1, 1, 0, 1, 0, 1, 1, 0},
					err:  nil,
				},
				{
					// before read: 110
					n:    2,
					bits: Stream{1, 1},
					err:  nil,
				},
				{
					// before read: 0
					n:    2,
					bits: Stream{0},
					err:  io.ErrUnexpectedEOF,
				},
			},
		},
	}

	for i, g := range golden {
		br := NewReader(g.r)
		for _, read := range g.reads {
			bits, err := br.Read(read.n)
			if err != read.err {
				t.Errorf("i %d: err; expected %q, got %q.", i, read.err, err)
			}
			if !Equal(bits, read.bits) {
				t.Errorf("i %d: bits; expected %q, got %q.", i, read.bits, bits)
			}
		}
	}
}

type testSeek struct {
	r     io.ReadSeeker
	reads []testReadParams
	seeks []testSeekParams
}

type testSeekParams struct {
	bitOff int64
	whence int
	bitRet int64
	err    error
}

// TestSeek alternates between calling Read and Seek, and verifies the return
// values.
func TestSeek(t *testing.T) {
	golden := []testSeek{
		// i = 0
		{
			// 10110111 01111011 11101111 11011111 11011111 11101111 11111011 11111111
			r: bytes.NewReader([]byte{0xB7, 0x7B, 0xEF, 0xDF, 0xDF, 0xEF, 0xFB, 0xFF}),
			reads: []testReadParams{
				{
					// offset before read: 0
					// before read: 10110111 01111011 11101111 11011111 11011111 11101111 11111011 11111111
					n:    1,
					bits: Stream{1},
					err:  nil,
				},
				{
					// offset before read: 4
					// before read: 0111 01111011 11101111 11011111 11011111 11101111 11111011 11111111
					n:    5,
					bits: Stream{0, 1, 1, 1, 0},
					err:  nil,
				},
				{
					// offset before read: 8
					// before read: 1111 11101111 11111011 11111111
					n:    28,
					bits: Stream{0, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 0, 1},
					err:  nil,
				},
				{
					// offset before read: 63
					// before read: 1
					n:    2,
					bits: Stream{1},
					err:  io.ErrUnexpectedEOF,
				},
			},
			seeks: []testSeekParams{
				{
					// offset before seek: 1
					bitOff: 3,
					whence: SeekCur,
					bitRet: 4,
					err:    nil,
				},
				{
					// offset before seek: 9
					bitOff: -1,
					whence: SeekCur,
					bitRet: 8,
					err:    nil,
				},
				{
					// offset before seek: 36
					bitOff: 27,
					whence: SeekCur,
					bitRet: 63,
					err:    nil,
				},
			},
		},

		// i = 1
		{
			// 10000010 10111000 00011010 10000011 10000001 10001010 10000011 00011010
			r: bytes.NewReader([]byte{0x82, 0xB8, 0x1A, 0x83, 0x81, 0x8A, 0x83, 0x1A}),
			reads: []testReadParams{
				{
					// offset before read: 0
					// before read: 10000010 10111000 00011010 10000011 10000001 10001010 10000011 00011010
					n:    0,
					bits: Stream{},
					err:  nil,
				},
				{
					// offset before read: 54
					// before read: 11 00011010
					n:    10,
					bits: Stream{1, 1, 0, 0, 0, 1, 1, 0, 1, 0},
					err:  nil,
				},
				{
					// offset before read: 7
					// before read: 0 10111000 00011010 10000011 10000001 10001010 10000011 00011010
					n:    5,
					bits: Stream{0, 1, 0, 1, 1},
					err:  nil,
				},
			},
			seeks: []testSeekParams{
				{
					// offset before seek: 0
					bitOff: -10,
					whence: SeekEnd,
					bitRet: 54,
					err:    nil,
				},
				{
					// offset before seek: 64
					bitOff: 7,
					whence: SeekSet,
					bitRet: 7,
					err:    nil,
				},
				{
					// offset before seek: 12
					bitOff: -13,
					whence: SeekCur,
					bitRet: 0,
					err:    errors.New("Seek: negative offset (-1)."),
				},
			},
		},

		// i = 2
		{
			// 10000010 10111000 00011010 10000011 10000001 10001010 10000011 00011010
			r:     bytes.NewReader([]byte{0x82, 0xB8, 0x1A, 0x83, 0x81, 0x8A, 0x83, 0x1A}),
			reads: []testReadParams{},
			seeks: []testSeekParams{
				{
					// offset before seek: 0
					bitOff: 7,
					whence: SeekCur,
					bitRet: 7,
					err:    nil,
				},
				{
					// offset before seek: 7
					bitOff: -15,
					whence: SeekCur,
					bitRet: 0,
					err:    errors.New("Seek: negative offset (-8)."),
				},
			},
		},

		// i = 3
		{
			// 10000010 10111000 00011010 10000011 10000001 10001010 10000011 00011010
			r:     bytes.NewReader([]byte{0x82, 0xB8, 0x1A, 0x83, 0x81, 0x8A, 0x83, 0x1A}),
			reads: []testReadParams{},
			seeks: []testSeekParams{
				{
					// offset before seek: 0
					bitOff: -1,
					whence: SeekEnd,
					bitRet: 63,
					err:    nil,
				},
				{
					// offset before seek: 63
					bitOff: 1,
					whence: SeekCur,
					bitRet: 64,
					err:    nil,
				},
				{
					// offset before seek: 64
					bitOff: 5,
					whence: SeekEnd,
					bitRet: 0,
					err:    errors.New("Seek: offset out of range; max 64, got 69."),
				},
			},
		},
	}

	for i, g := range golden {
		br := NewReader(g.r)
		max := len(g.reads)
		if max < len(g.seeks) {
			max = len(g.seeks)
		}
		for j := 0; j < max; j++ {
			if len(g.reads) > j {
				read := g.reads[j]
				bits, err := br.Read(read.n)
				if err != read.err {
					t.Errorf("i %d, j %d: err; expected %q, got %q.", i, j, read.err, err)
				}
				if !Equal(bits, read.bits) {
					t.Errorf("i %d, j %d: bits; expected %q, got %q.", i, j, read.bits, bits)
				}
			}
			if len(g.seeks) > j {
				seek := g.seeks[j]
				bitRet, err := br.Seek(seek.bitOff, seek.whence)
				e1 := err
				e2 := seek.err
				if e1 != nil && e2 != nil {
					if e1.Error() != e2.Error() {
						t.Errorf("i %d, j %d: err; expected %q, got %q.", i, j, seek.err, err)
					}
				} else if e1 != nil || e2 != nil {
					t.Errorf("i %d, j %d: err; expected %q, got %q.", i, j, seek.err, err)
				}
				if bitRet != seek.bitRet {
					t.Errorf("i %d, j %d: bitRet; expected %v, got %v.", i, j, seek.bitRet, bitRet)
				}
			}
		}
	}
}
