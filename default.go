package log

import (
	"fmt"
	"io"
	gl "log"
	"os"
	"time"
)

var defaultLogger = &Log{
	w:    ref(os.Stdout),
	line: 1,
}

// To sets the writer for the default logger
func To(w io.Writer) { defaultLogger.To(w) }

// ToFile sets the default logger to log to a file. If no path is given, it will
// creat a file in the working directory with the same name as the program
// followed by ".log". If more than one path is given, everything beyond the
// first is ingored.
func ToFile(path ...string) error {
	f, err := open(path...)
	if err != nil {
		return err
	}
	defaultLogger.To(f)
	return nil
}

// Child method for the default logger
func Child(data ...interface{}) *Log { return defaultLogger.Child(data...) }

// SetDebug for the default logger
func SetDebug(debug bool) { defaultLogger.SetDebug(debug) }

// GetDebug for the default logger
func GetDebug() bool { return defaultLogger.GetDebug() }

// Debug for the default logger
func Debug(data ...interface{}) { defaultLogger.Debug(data...) }

// Info writes data to the default logger
func Info(data ...interface{}) { defaultLogger.Info(data...) }

// Error method for the default logger
func Error(err error) bool { return defaultLogger.Error(err) }

// Panic method for the defulat logger
func Panic(err error) { defaultLogger.Panic(err) }

// Fatal method for the defulat logger
func Fatal(err error) { defaultLogger.Fatal(err) }

type glWrapper struct{}

func (g *glWrapper) Write(b []byte) (int, error) {
	w := defaultLogger.w
	w.Lock()
	defer w.Unlock()
	fmt.Fprint(w.Writer, Formatter(time.Now()))
	return w.Writer.Write(b)
}

var glw = &glWrapper{}

// Go wraps log from the Go standard library so it will write the default log.
func Go() {
	gl.SetPrefix(" GOLOG ")
	gl.SetFlags(0)
	gl.SetOutput(glw)
}
