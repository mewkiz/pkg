// Package bytesutil implements some bytes utility functions.
package bytesutil

import "bytes"

// IndexAfter returns the index directly after the first instance of sep in s,
// or -1 if sep is not present in s.
func IndexAfter(s, sep []byte) int {
   pos := bytes.Index(s, sep)
   if pos == -1 {
      return -1
   }
   return pos + len(sep)
}
