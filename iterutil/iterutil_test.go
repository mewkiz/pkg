package iterutil_test

import (
	"reflect"
	"testing"

	"github.com/mewkiz/pkg/iterutil"
)

func TestOrdered(t *testing.T) {
	golden := []struct {
		input map[int]string
		want  []int
	}{
		{
			input: map[int]string{
				11: "eleven",
				10: "ten",
				9:  "nine",
				5:  "five",
				12: "twelve",
				6:  "six",
				8:  "eight",
				2:  "two",
				13: "thirteen",
				4:  "four",
				3:  "three",
				1:  "one",
				7:  "seven",
			},
			want: []int{
				1,
				2,
				3,
				4,
				5,
				6,
				7,
				8,
				9,
				10,
				11,
				12,
				13,
			},
		},
	}

	for _, g := range golden {
		var got []int
		for key := range iterutil.Ordered(g.input) {
			got = append(got, key)
		}
		if !reflect.DeepEqual(g.want, got) {
			t.Errorf("key order mismatch; expected %v, got %v", g.want, got)
			continue
		}
	}
}
