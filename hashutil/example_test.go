package hashutil_test

import (
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os"
	"strings"

	"github.com/mewkiz/pkg/hashutil"
)

func ExampleHashReader() {
	r := strings.NewReader("The quick brown fox jumps over the lazy dog.\n")
	h := crc32.NewIEEE()
	hr := hashutil.NewHashReader(r, h)
	_, err := io.Copy(os.Stdout, hr)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("checksum: 0x%08X", hr.Sum32())
	// Output:
	// The quick brown fox jumps over the lazy dog.
	// checksum: 0xEB50CC6A
}
