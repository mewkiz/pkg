// Package bufioutil implements some bufio utility functions.
package bufioutil

import "bufio"
import "io"

// ReadLine returns a single line, not including the end-of-line bytes.
func ReadLine(br *bufio.Reader) (line string, err error) {
   for {
      buf, isPrefix, err := br.ReadLine()
      if err != nil {
         return "", err
      }
      line += string(buf)
      if !isPrefix {
         break
      }
   }
   return line, nil
}

// ReadLines returns all lines, not including the end-of-line bytes.
func ReadLines(br *bufio.Reader) (lines []string, err error) {
   for {
      line, err := ReadLine(br)
      if err != nil {
         if err == io.EOF {
            break
         }
         return nil, err
      }
      lines = append(lines, line)
   }
   return lines, nil
}
