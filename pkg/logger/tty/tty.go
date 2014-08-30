//
// Logger implementation for non-interactive shells and/or verbose mode.
//
package tty_logger

import "time"
import "fmt"
import "io"
import "os"

type Logger struct {
	start time.Time
	w     io.Writer
}

// New logger with the given writer.
func New(w io.Writer) *Logger {
	return &Logger{
		start: time.Now(),
		w:     w,
	}
}

// Log skipped commit.
func (l *Logger) Skip(arg, commit string) {
	fmt.Fprintf(l.w, "\033[0m   %s %s\033[0m\n", shorten(commit), arg)
}

// Log running of commit.
func (l *Logger) Start(arg, commit string) {
	fmt.Fprintf(l.w, "\033[36m   %s\033[0m %s \n", shorten(commit), arg)
}

// Log successful commit.
func (l *Logger) Success(arg, commit string) {
}

// Log error.
func (l *Logger) Error(err error, arg, commit string) {
	fmt.Fprintf(l.w, "\n\033[31m   %s %s: %s\033[0m\n", arg, commit, err)
}

// Log line.
func (l *Logger) Log(line string) {
	fmt.Fprintf(l.w, "\n\033[90m   %s\033[0m\n", line)
}

// Log end of provisioning.
func (l *Logger) End() {
	fmt.Fprintf(l.w, "\033[32m   completed\033[0m in %s\n", time.Since(l.start))
}

// Stdout implementation.
func (l *Logger) Stdout() io.Writer {
	return os.Stdout
}

// Stderr implementation.
func (l *Logger) Stderr() io.Writer {
	return os.Stderr
}

// return shortened hash.
func shorten(commit string) string {
	return commit[:8]
}
