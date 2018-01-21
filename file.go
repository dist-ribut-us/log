package log

import (
	"os"
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
