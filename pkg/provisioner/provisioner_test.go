package provisioner

import "github.com/visionmedia/stack/pkg/logger/tty"
import "github.com/bmizerany/assert"
import "testing"
import "bytes"
import "os"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Test run with comments and newlines
func TestRun(t *testing.T) {
	os.Remove("/tmp/commits.log")

	r := bytes.NewBuffer(nil)

	r.WriteString("echo foo\n")
	r.WriteString("echo bar\n")
	r.WriteString("\n")
	r.WriteString("  \n")
	r.WriteString("  \r\n")
	r.WriteString("  \r")
	r.WriteString("# RUN echo bar\n")
	r.WriteString("# something here\n")
	r.WriteString("echo baz\n")

	p := New(r, "/tmp/commits.log")
	p.Log = tty_logger.New(os.Stdout)

	check(p.Run())

	assert.Equal(t, 3, p.commits.Length())
}

// Test with canonical commands.
func TestRunCanonical(t *testing.T) {
	os.Remove("/tmp/commits.log")

	r := bytes.NewBuffer(nil)

	r.WriteString("RUN echo foo\n")
	r.WriteString("RUN echo bar\n")
	r.WriteString("\n")
	r.WriteString("  \n")
	r.WriteString("  \r\n")
	r.WriteString("  \r")
	r.WriteString("RUN echo bar\n")
	r.WriteString("LOG something here\n")
	r.WriteString("RUN echo baz\n")

	p := New(r, "/tmp/commits.log")
	p.Log = tty_logger.New(os.Stdout)

	check(p.Run())

	assert.Equal(t, 3, p.commits.Length())
}
