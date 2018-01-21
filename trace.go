package log

import (
	"fmt"
	"runtime"
	"strings"
)

// Line returns the invoking line
func Line() Val {
	// pop 2 to remove calls to Line and Trace
	return Trace(2, 1)
}

// CallLine returns the a single line above steps above in the stack trace
func CallLine(above int) Val {
	// pop 2 to remove calls to Line and Trace
	return Trace(above+2, 1)
}

// Trace calls TraceDepth using PathDepth for the depth.
func Trace(start, end int) Val {
	return TraceDepth(start, end, PathDepth)
}

// TraceDepth returns a formatted slice of the stack trace. The start arg
// determines how close to the calling line to return and end determines how
// close to the top of the stack trace to return.
//
// If start is positive, that many lines from the beginning will be removed. If
// start is negative, that is how many lines will be returned (and end will be
// the anchor).
//
// If end is negative, that is how far from the top of the stack trace end will
// be. If end is positive, that is how many lines will be returned (and start
// will be the anchor).
//
// If start is negative and end is positive, then I don't know what happens!
func TraceDepth(start, end, pathDepth int) Val {
	if end == 0 || (start < 0 && end >= 0) {
		return nil
	}

	if start >= 0 && end == 1 {
		str, _ := caller(start+2, pathDepth)
		return Lbl(str)
	}

	var lines []string
	linesStart := start

	if linesStart < 0 {
		linesStart = 0
	}

	for i := linesStart + 2; true; i++ {
		line, ok := caller(i, pathDepth)
		if !ok {
			break
		}
		lines = append(lines, line)
	}

	if start < 0 {
		start += len(lines)
		if start < 0 {
			start = 0
		}
		lines = lines[start:]
	}

	if end < 0 {
		end += len(lines)
		if end <= 0 {
			return nil
		}
	}
	if end < len(lines) {
		lines = lines[:end]
	}

	return Lbl("\n\t" + strings.Join(lines, "\n\t"))
}

func caller(i, pd int) (string, bool) {
	ptr, file, ln, ok := runtime.Caller(i)
	if !ok {
		return "", false
	}
	return fmt.Sprintf("%s:%d %s", getFileName(file, pd), ln, getFuncName(ptr)), true
}

func getFileName(file string, pd int) string {
	if pd < 0 {
		return file
	}
	for i := len(file) - 1; i >= 0; i-- {
		if file[i] == '/' {
			if pd == 0 {
				return file[i+1:]
			}
			pd--
		}
	}
	return file
}

func getFuncName(ptr uintptr) string {
	f := runtime.FuncForPC(ptr)
	if f == nil {
		return ""
	}
	return f.Name()
}
