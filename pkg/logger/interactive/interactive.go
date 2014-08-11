//
// Logger implementation for interactive shell sessions.
//
package interactive_logger

import "github.com/visionmedia/go-spin"
import "time"
import "fmt"
import "io"

type Logger struct {
	start time.Time
	w     io.Writer
	q     chan bool
}

// New logger with the given writer.
func New(w io.Writer) *Logger {
	return &Logger{
		start: time.Now(),
		q:     make(chan bool),
		w:     w,
	}
}

// Log skipped commit.
func (l *Logger) Skip(arg, commit string) {
	fmt.Fprintf(l.w, "\033[0m   %s %s\033[0m\n", shorten(commit), arg)
}

// Log running of commit.
func (l *Logger) Start(arg, commit string) {
	s := spin.New()

	go func() {
		for {
			select {
			case <-l.q:
				return
			case <-time.Tick(100 * time.Millisecond):
				fmt.Fprintf(l.w, "\r\033[36m %s %s\033[0m %s ", s.Next(), shorten(commit), arg)
				s.Next()
			}
		}
	}()

	fmt.Fprintf(l.w, "\033[36m   %s\033[0m %s ", shorten(commit), arg)
}

// Log successful commit.
func (l *Logger) Success(arg, commit string) {
	l.q <- true
	fmt.Fprintf(l.w, "\r\033[32m   %s\033[0m %s\n", shorten(commit), arg)
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

// return shortened hash.
func shorten(commit string) string {
	return commit[:8]
}
