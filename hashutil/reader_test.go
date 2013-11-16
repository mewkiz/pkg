package hashutil_test

import (
	"bytes"
	"hash"
	"hash/crc32"
	"hash/crc64"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/mewkiz/pkg/hashutil"
	"github.com/mewkiz/pkg/hashutil/crc16"
	"github.com/mewkiz/pkg/hashutil/crc8"
)

type test struct {
	crc8  uint8
	crc16 uint16
	crc32 uint32
	crc64 uint64
	in    string
}

var golden = []test{
	{0x00, 0x0000, 0x00000000, 0x0000000000000000, ""},
	{0x20, 0x8145, 0xE8B7BE43, 0x3420000000000000, "a"},
	{0xC9, 0xC749, 0x9E83486D, 0x36C4200000000000, "ab"},
	{0x5F, 0xCADB, 0x352441C2, 0x3776C42000000000, "abc"},
	{0xA1, 0x58E7, 0xED82CD11, 0x336776C420000000, "abcd"},
	{0x52, 0x678D, 0x8587D865, 0x32D36776C4200000, "abcde"},
	{0x8C, 0x0D05, 0x4B8E39EF, 0x3002D36776C42000, "abcdef"},
	{0x9F, 0x047C, 0x312A6AA6, 0x31B002D36776C420, "abcdefg"},
	{0xCB, 0x7D68, 0xAEEF2A50, 0x0E21B002D36776C4, "abcdefgh"},
	{0x67, 0x6878, 0x8DA988AF, 0x8B6E21B002D36776, "abcdefghi"},
	{0x23, 0xF80F, 0x3981703A, 0x7F5B6E21B002D367, "abcdefghij"},
	{0x56, 0x0F8E, 0x6B9CDFE7, 0x8EC0E7C835BF9CDF, "Discard medicine more than two years old."},
	{0x6B, 0xE149, 0xC90EF73F, 0xC7DB1759E2BE5AB4, "He who has a shady past knows that nice guys finish last."},
	{0x70, 0x02B7, 0xB902341F, 0xFBF9D9603A6FA020, "I wouldn't marry him with a ten foot pole."},
	{0x8F, 0x7F6A, 0x042080E8, 0xEAFC4211A6DAA0EF, "Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave"},
	{0x48, 0x28BD, 0x154C6D11, 0x3E05B21C7A4DC4DA, "The days of the digital watch are numbered.  -Tom Stoppard"},
	{0x5E, 0x7C55, 0x4C418325, 0x5255866AD6EF28A6, "Nepal premier won't resign."},
	{0x3C, 0xC92B, 0x33955150, 0x8A79895BE1E9C361, "For every action there is an equal and opposite government program."},
	{0xA8, 0x3E41, 0x26216A4B, 0x8878963A649D4916, "His money is twice tainted: 'taint yours and 'taint mine."},
	{0x46, 0xDA56, 0x1ABBE45E, 0xA7B9D53EA87EB82F, "There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977"},
	{0xC7, 0x7F66, 0xC89A94F7, 0xDB6805C0966A2F9C, "It's a tiny change to the code and not completely disgusting. - Bob Manchek"},
	{0x31, 0x2A00, 0xAB3ABE14, 0xF3553C65DACDADD2, "size:  a.out:  bad magic"},
	{0xB6, 0x25B2, 0xBAB102B6, 0x9D5E034087A676B9, "The major problem is with sendmail.  -Mark Horton"},
	{0x7D, 0xBD71, 0x999149D7, 0xA6DB2D7F8DA96417, "Give me a rock, paper and scissors and I will move the world.  CCFestoon"},
	{0xDC, 0x8596, 0x6D52A33C, 0x325E00CD2FE819F9, "If the enemy is within range, then so are you."},
	{0x13, 0x74A2, 0x90631E8D, 0x88C6600CE58AE4C6, "It's well we cannot hear the screams/That we create in others' dreams."},
	{0x96, 0x0D73, 0x78309130, 0x28C4A3F3B769E078, "You remind me of a TV show, but that's all right: I watch it anyway."},
	{0x96, 0xEE65, 0x7D0A377F, 0xA698A34C9D9F1DCA, "C is as portable as Stonehedge!!"},
	{0x3C, 0xA94E, 0x8C79FD79, 0xF6C1E2A8C26C5CFC, "Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley"},
	{0xEE, 0x0B98, 0xA20B7167, 0x0D402559DFE9B70C, "The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule"},
	{0x33, 0xF560, 0x8E0BB443, 0xDB6EFFF26AA94946, "How can you write a big system without C++?  -Paul Glick"},
}

func TestHashReaderCrc8(t *testing.T) {
	for i, g := range golden {
		hr := hashutil.NewHashReader(strings.NewReader(g.in), crc8.NewATM())
		buf, err := ioutil.ReadAll(hr)
		if err != nil {
			t.Errorf("i=%d: error during read; %v", i, err)
			continue
		}
		s := string(buf)
		if g.in != s {
			t.Errorf("i=%d: expected %q, got %q", g.in, s)
			continue
		}
		got := hr.Sum8()
		if got != g.crc8 {
			t.Errorf("i=%d: crc8(%q); expected 0x%02X, want 0x%02X", i, g.in, got, g.crc8)
		}
	}
}

func TestHashReaderCrc16(t *testing.T) {
	for i, g := range golden {
		hr := hashutil.NewHashReader(strings.NewReader(g.in), crc16.NewIBM())
		buf, err := ioutil.ReadAll(hr)
		if err != nil {
			t.Errorf("i=%d: error during read; %v", i, err)
			continue
		}
		s := string(buf)
		if g.in != s {
			t.Errorf("i=%d: expected %q, got %q", g.in, s)
			continue
		}
		got := hr.Sum16()
		if got != g.crc16 {
			t.Errorf("i=%d: crc16(%q); expected 0x%04X, want 0x%04X", i, g.in, got, g.crc16)
		}
	}
}

func TestHashReaderCrc32(t *testing.T) {
	for i, g := range golden {
		hr := hashutil.NewHashReader(strings.NewReader(g.in), crc32.NewIEEE())
		buf, err := ioutil.ReadAll(hr)
		if err != nil {
			t.Errorf("i=%d: error during read; %v", i, err)
			continue
		}
		s := string(buf)
		if g.in != s {
			t.Errorf("i=%d: expected %q, got %q", g.in, s)
			continue
		}
		got := hr.Sum32()
		if got != g.crc32 {
			t.Errorf("i=%d: crc32(%q); expected 0x%04X, want 0x%04X", i, g.in, got, g.crc32)
		}
	}
}

var table64 = crc64.MakeTable(crc64.ISO)

func TestHashReaderCrc64(t *testing.T) {
	for i, g := range golden {
		hr := hashutil.NewHashReader(strings.NewReader(g.in), crc64.New(table64))
		buf, err := ioutil.ReadAll(hr)
		if err != nil {
			t.Errorf("i=%d: error during read; %v", i, err)
			continue
		}
		s := string(buf)
		if g.in != s {
			t.Errorf("i=%d: expected %q, got %q", g.in, s)
			continue
		}
		got := hr.Sum64()
		if got != g.crc64 {
			t.Errorf("i=%d: crc64(%q); expected 0x%04X, want 0x%04X", i, g.in, got, g.crc64)
		}
	}
}

func BenchmarkHashReaderCrc8_1K(b *testing.B) {
	h := crc8.NewATM()
	benchmarkHashReader(b, h, 1024)
}

func BenchmarkHashReaderCrc8_2K(b *testing.B) {
	h := crc8.NewATM()
	benchmarkHashReader(b, h, 2*1024)
}

func BenchmarkHashReaderCrc8_4K(b *testing.B) {
	h := crc8.NewATM()
	benchmarkHashReader(b, h, 4*1024)
}

func BenchmarkHashReaderCrc8_8K(b *testing.B) {
	h := crc8.NewATM()
	benchmarkHashReader(b, h, 8*1024)
}

func BenchmarkHashReaderCrc8_16K(b *testing.B) {
	h := crc8.NewATM()
	benchmarkHashReader(b, h, 16*1024)
}

func BenchmarkHashReaderCrc16_1K(b *testing.B) {
	h := crc16.NewIBM()
	benchmarkHashReader(b, h, 1024)
}

func BenchmarkHashReaderCrc16_2K(b *testing.B) {
	h := crc16.NewIBM()
	benchmarkHashReader(b, h, 2*1024)
}

func BenchmarkHashReaderCrc16_4K(b *testing.B) {
	h := crc16.NewIBM()
	benchmarkHashReader(b, h, 4*1024)
}

func BenchmarkHashReaderCrc16_8K(b *testing.B) {
	h := crc16.NewIBM()
	benchmarkHashReader(b, h, 8*1024)
}

func BenchmarkHashReaderCrc16_16K(b *testing.B) {
	h := crc16.NewIBM()
	benchmarkHashReader(b, h, 16*1024)
}

func BenchmarkHashReaderCrc32_1K(b *testing.B) {
	h := crc32.NewIEEE()
	benchmarkHashReader(b, h, 1024)
}

func BenchmarkHashReaderCrc32_2K(b *testing.B) {
	h := crc32.NewIEEE()
	benchmarkHashReader(b, h, 2*1024)
}

func BenchmarkHashReaderCrc32_4K(b *testing.B) {
	h := crc32.NewIEEE()
	benchmarkHashReader(b, h, 4*1024)
}

func BenchmarkHashReaderCrc32_8K(b *testing.B) {
	h := crc32.NewIEEE()
	benchmarkHashReader(b, h, 8*1024)
}

func BenchmarkHashReaderCrc32_16K(b *testing.B) {
	h := crc32.NewIEEE()
	benchmarkHashReader(b, h, 16*1024)
}

func BenchmarkHashReaderCrc64_1K(b *testing.B) {
	h := crc64.New(table64)
	benchmarkHashReader(b, h, 1024)
}

func BenchmarkHashReaderCrc64_2K(b *testing.B) {
	h := crc64.New(table64)
	benchmarkHashReader(b, h, 2*1024)
}

func BenchmarkHashReaderCrc64_4K(b *testing.B) {
	h := crc64.New(table64)
	benchmarkHashReader(b, h, 4*1024)
}

func BenchmarkHashReaderCrc64_8K(b *testing.B) {
	h := crc64.New(table64)
	benchmarkHashReader(b, h, 8*1024)
}

func BenchmarkHashReaderCrc64_16K(b *testing.B) {
	h := crc64.New(table64)
	benchmarkHashReader(b, h, 16*1024)
}

func benchmarkHashReader(b *testing.B, h hash.Hash, count int64) {
	b.SetBytes(count)
	data := make([]byte, count)
	for i := range data {
		data[i] = byte(i)
	}
	in := make([]byte, 0, h.Size())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Reset()
		hr := hashutil.NewHashReader(bytes.NewReader(data), h)
		io.Copy(ioutil.Discard, hr)
		h.Sum(in)
	}
}
