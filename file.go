package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// File opens a file as a log. If no path is given, it will creat a file in the
// working directory with the same name as the program followed by ".log". If
// more than one path is given, everything beyond the first is ingored.
func File(path ...string) (*Log, error) {
	f, err := open(path...)
	if err != nil {
		return nil, err
	}
	return New(f), nil
}

const (
	// Truncate is a reference to os.O_TRUNC
	Truncate = os.O_TRUNC

	// Append is a reference to os.O_APPEND
	Append = os.O_APPEND
)

// Contents determines how the existing contents of the log will be treated. It
// should be set to either Truncate or Append.
var Contents = Append

func open(path ...string) (*os.File, error) {
	var p string
	if len(path) > 0 {
		p = path[0]
	} else {
		path = strings.Split(os.Args[0], "/")
		p = strings.Split(path[len(path)-1], ".")[0] + ".log"
	}
	return os.OpenFile(p, Contents|os.O_CREATE|os.O_WRONLY, 0666)
}

// PathDepth sets how many directories to include when showing file names in
// calls to Line
var PathDepth = 1

// Trim removes the bottom layers from the stack trace when calling Line. These
// are always the runtime and so not useful information.
var Trim uint = 2

// Line can get either a single line or a stack trace. If no arguments are given
// it will return the line from which it was invoked. If a positive int is given
// it will return the line that many steps up the stack. If a negative int is
// given, it will produce a stack trace starting that many steps up the stack.
//
// It also optionally takes a second int to set how many directory levels should
// be included when showing file names. If this is negative, the entire path
// will be included. If it is omitted, PathDepth will be used.
//
// A third argument will set the number of layers to trim from the bottom of the
// stack. If it is negative it will print exactly that many layers. If it is
// omitted, Trim will be used.
//
// tl;dr Just include log.Line() in a log call to include the file and line
// number or put log.Line(-1) as the last argument a log call to include a stack
// trace.
func Line(ints ...int) Val {
	l := len(ints)
	i := 0
	pd := PathDepth
	trim := int(Trim)
	if l > 0 {
		i = ints[0]
		if l > 1 {
			pd = ints[1]
			if i < 0 && l > 2 {
				trim = ints[2]
			}
		}
	}
	if i < 0 { // return trace
		strs := []string{""}
		for i = 1 - i; true; i++ {
			str, ok := caller(i, pd)
			if !ok {
				break
			}
			strs = append(strs, str)
		}
		if trim >= 0 {
			l = len(strs) - trim - 1
		} else {
			l = len(strs) - 1
			trim = -trim + 1
			if trim < l {
				l = trim
			}
		}
		if l < 0 {
			l = 0
		}
		return Lbl(strings.Join(strs[:l], "\n  "))
	} // else return single line
	str, _ := caller(i+2, pd)
	return Lbl(str)
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
