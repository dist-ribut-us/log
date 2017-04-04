package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// Lbl provides a way to add strings that will not be wrapped in quotes
type Lbl string

// LogVal returns the underlying string
func (l Lbl) LogVal() string { return string(l) }

// KV write a key value pair as "key=value".
type KV struct {
	K string
	V interface{}
}

// LogVal sets how KV will be displayed in the logs
func (kv KV) LogVal() string {
	return fmt.Sprintf("%s=%s", kv.K, Formatter(kv.V))
}

// Val takes priority in Format.
type Val interface {
	LogVal() string
}

// TimeFormat to format time in the log
var TimeFormat = "2006-01-02_15:04:05.00"

// Format takes a value and converts it to a string. It makes sure any string
// values are wrapped in quotes. Anything with a method of Log()string
func Format(value interface{}) string {
	switch t := value.(type) {
	case Val:
		return t.LogVal()
	case time.Time:
		return t.UTC().Format(TimeFormat)
	case string, error, fmt.Stringer:
		return fmt.Sprintf("\"%s\"", t)
	}
	return fmt.Sprint(value)
}

// Formatter is what the log uses to format values to strings. It defaults to
// Format, but can be replaced.
var Formatter = Format

// Log writes data to a log. Each entry is one line
type Log struct {
	w     *wRef
	data  string
	debug bool
	line  int
	trace []int
}

// wRef is a reference to the io.Writer and it's mutex. This allows all the
// children of one log to share the same reference. If it is changed in one
// place, it will change everywhere.
type wRef struct {
	io.Writer
	*sync.Mutex
}

func ref(w io.Writer) *wRef {
	return &wRef{w, &sync.Mutex{}}
}

func (r *wRef) Write(b []byte) (int, error) {
	if r == nil || r.Writer == nil {
		return 0, nil
	}
	r.Lock()
	defer r.Unlock()
	return r.Writer.Write(b)
}

// New creates a new log from a writer.
func New(w io.Writer) *Log {
	return &Log{
		w: ref(w),
	}
}

// To changes the writer for a log
func (l *Log) To(w io.Writer) {
	if l == nil {
		l = &Log{
			w: ref(w),
		}
		return
	}
	l.w.Lock()
	defer l.w.Unlock()
	l.w.Writer = w
}

// Mute a log by setting the writer to nil
func (l *Log) Mute() { l.To(nil) }

// SetTrace to include line in log messages
func (l *Log) SetTrace(ints ...int) {
	if len(ints) == 0 {
		ints = []int{l.line}
	}

	if ints[0] >= 0 {
		ints[0] += 2
	} else {
		ints[0] -= 2
	}
	l.trace = ints
}

// NoTrace disables SetTrace
func (l *Log) NoTrace() {
	l.trace = nil
}

// Child creates a new log that uses the same writer. All data passed in will
// be written on every line written to this log.
func (l *Log) Child(data ...interface{}) *Log {
	if l == nil {
		return nil
	}
	var strs []string
	o := 1 //offset
	if l.data == "" {
		o = 0
		strs = make([]string, len(data))
	} else {
		strs = make([]string, len(data)+1)
		strs[0] = l.data
	}
	for i, d := range data {
		strs[i+o] = Formatter(d)
	}
	return &Log{
		w:     l.w,
		data:  strings.Join(strs, " "),
		debug: l.debug,
		trace: l.trace,
	}
}

// SetDebug sets the debug bool
func (l *Log) SetDebug(debug bool) {
	if l == nil {
		return
	}
	l.debug = debug
}

// GetDebug gets the debug bool
func (l *Log) GetDebug() bool {
	return l != nil && l.debug
}

// Debug will write data to the log only if debug is enabled
func (l *Log) Debug(data ...interface{}) {
	if l == nil || !l.debug {
		return
	}
	l.write("DEBUG", data...)
}

func (l *Log) write(flag string, data ...interface{}) {
	if l == nil || l.w.Writer == nil {
		return
	}

	// find length and offset of strings
	o := 2
	ldata := -1
	trace := -1
	ln := len(data)
	if len(l.trace) > 0 && (flag != "ERROR" && flag != "DEBUG") {
		if l.trace[0] >= 0 {
			trace = o
			o++
		} else {
			trace = -2
			ln++
		}
	}
	if l.data != "" {
		ldata = o
		o++
	}

	strs := make([]string, ln+o)

	// populate strings
	strs[0] = Formatter(time.Now())
	strs[1] = flag
	if trace == -2 {
		strs[ln+o-1] = Line(l.trace...).LogVal()
	} else if trace > -1 {
		strs[trace] = Line(l.trace...).LogVal()
	}
	if ldata > -1 {
		strs[ldata] = l.data
	}
	for i, d := range data {
		strs[i+o] = Formatter(d)
	}
	fmt.Fprintln(l.w, strings.Join(strs, " "))
}

// Info writes data to the log with the INFO flag
func (l *Log) Info(data ...interface{}) { l.write("INFO", data...) }

// Error takes an error and if it is not nil writes it to the log. It returns
// a bool indicating if there was an error.
func (l *Log) Error(err error) bool {
	if err == nil {
		return false
	}
	l.write("ERROR", Line(1+l.line), err)
	return true
}

// Panic takes an error and if it is not nil, writes it to the log then panics.
// Panic should only be called from a main package.
func (l *Log) Panic(err error) {
	if err == nil {
		return
	}
	l.write("PANIC", Line(1+l.line), err)
	l.Close()
	panic(err)
}

// Fatal writes data to the log with the FATAL flag and call os.Exit(1). Fatal
// should only be called from a main package.
func (l *Log) Fatal(data ...interface{}) {
	l.write("FATAL", data...)
	l.Close()
	os.Exit(1)
}

// Close will close the writer if the underlying type has a Close method. It
// will only return an error if one is generated while closing the writer. If
// the writer is not actually a WriterCloser, no error will be returned.
func (l *Log) Close() error {
	if l == nil {
		return nil
	}
	if wc, ok := l.w.Writer.(io.WriteCloser); ok {
		return wc.Close()
	}
	return nil
}
