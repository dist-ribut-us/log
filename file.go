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

// Line attempts to return the file and line number i positions up the call
// stack. It returns only the file name, not the directory. In the case that it
// cannot fetch the information, an empty string is returned.
func Line(i int) Val {
	_, file, l, ok := runtime.Caller(i + 1)
	if !ok {
		return line("")
	}
	for i := len(file) - 1; i >= 0; i-- {
		if file[i] == '/' {
			file = file[i+1:]
			break
		}
	}
	return line(fmt.Sprint(file, ":", l))
}

type line string

func (l line) LogVal() string { return string(l) }
