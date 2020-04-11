// Package stackutil provides access to stack traces.
package stackutil

import (
	"fmt"
	"runtime"
	"strings"
)

// StackTrace returns a stack trace of the current function caller stack.
func StackTrace() string {
	pc := make([]uintptr, 10)
	n := runtime.Callers(2, pc) // skip runtime.Callers and stackutil.StackTrace caller pcs.
	if n == 0 {
		return ""
	}
	pc = pc[:n]
	frames := runtime.CallersFrames(pc)
	buf := &strings.Builder{}
	for {
		frame, more := frames.Next()
		fmt.Fprintf(buf, "%s:%d\n", frame.Function, frame.Line)
		if !more {
			break
		}
	}
	return buf.String()
}
