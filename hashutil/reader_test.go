package hashutil_test

import (
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
	in    string
}

var golden = []test{
	{0x00, 0x0000, ""},
	{0x20, 0x8145, "a"},
	{0xC9, 0xC749, "ab"},
	{0x5F, 0xCADB, "abc"},
	{0xA1, 0x58E7, "abcd"},
	{0x52, 0x678D, "abcde"},
	{0x8C, 0x0D05, "abcdef"},
	{0x9F, 0x047C, "abcdefg"},
	{0xCB, 0x7D68, "abcdefgh"},
	{0x67, 0x6878, "abcdefghi"},
	{0x23, 0xF80F, "abcdefghij"},
	{0x56, 0x0F8E, "Discard medicine more than two years old."},
	{0x6B, 0xE149, "He who has a shady past knows that nice guys finish last."},
	{0x70, 0x02B7, "I wouldn't marry him with a ten foot pole."},
	{0x8F, 0x7F6A, "Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave"},
	{0x48, 0x28BD, "The days of the digital watch are numbered.  -Tom Stoppard"},
	{0x5E, 0x7C55, "Nepal premier won't resign."},
	{0x3C, 0xC92B, "For every action there is an equal and opposite government program."},
	{0xA8, 0x3E41, "His money is twice tainted: 'taint yours and 'taint mine."},
	{0x46, 0xDA56, "There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977"},
	{0xC7, 0x7F66, "It's a tiny change to the code and not completely disgusting. - Bob Manchek"},
	{0x31, 0x2A00, "size:  a.out:  bad magic"},
	{0xB6, 0x25B2, "The major problem is with sendmail.  -Mark Horton"},
	{0x7D, 0xBD71, "Give me a rock, paper and scissors and I will move the world.  CCFestoon"},
	{0xDC, 0x8596, "If the enemy is within range, then so are you."},
	{0x13, 0x74A2, "It's well we cannot hear the screams/That we create in others' dreams."},
	{0x96, 0x0D73, "You remind me of a TV show, but that's all right: I watch it anyway."},
	{0x96, 0xEE65, "C is as portable as Stonehedge!!"},
	{0x3C, 0xA94E, "Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley"},
	{0xEE, 0x0B98, "The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule"},
	{0x33, 0xF560, "How can you write a big system without C++?  -Paul Glick"},
}

func TestHashReaderCRC8(t *testing.T) {
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

func TestHashReaderCRC16(t *testing.T) {
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
