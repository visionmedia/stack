//
// Logger implementation for color haters.
//
package plain_logger

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
	fmt.Fprintf(l.w, "   %s %s\n", shorten(commit), arg)
}

// Log running of commit.
func (l *Logger) Start(arg, commit string) {
	fmt.Fprintf(l.w, "   %s %s \n\n", shorten(commit), arg)
}

// Log successful commit.
func (l *Logger) Success(arg, commit string) {
}

// Log error.
func (l *Logger) Error(err error, arg, commit string) {
	fmt.Fprintf(l.w, "\n   %s %s: %s\n", arg, commit, err)
}

// Log line.
func (l *Logger) Log(line string) {
	fmt.Fprintf(l.w, "\n   %s\n", line)
}

// Log end of provisioning.
func (l *Logger) End() {
	fmt.Fprintf(l.w, "   completed in %s\n", time.Since(l.start))
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
