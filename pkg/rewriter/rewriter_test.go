package rewriter

import "github.com/bmizerany/assert"
import "testing"
import "bytes"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func TestRun(t *testing.T) {
	in := bytes.NewBuffer(nil)

	in.WriteString("echo foo\n")
	in.WriteString("echo bar\n")
	in.WriteString("   \r\n")
	in.WriteString("\n")
	in.WriteString("echo baz")

	out, err := Rewrite(in)
	check(err)

	exp := `RUN echo foo
RUN echo bar
RUN echo baz
`

	assert.Equal(t, exp, out)
}

func TestLog(t *testing.T) {
	in := bytes.NewBuffer(nil)

	in.WriteString("# Install node\n")
	in.WriteString("LOG Install node\n")
	in.WriteString("echo foo\n")
	in.WriteString("echo bar\n")

	out, err := Rewrite(in)
	check(err)

	exp := `LOG Install node
LOG Install node
RUN echo foo
RUN echo bar
`

	assert.Equal(t, exp, out)
}

func TestInclude(t *testing.T) {
	in := bytes.NewBuffer(nil)

	in.WriteString(". fixtures/a.sh\n")
	in.WriteString("include fixtures/b.sh\n")

	out, err := Rewrite(in)
	check(err)

	exp := `RUN echo "from a"
RUN echo "from a again"
RUN echo "from b"
RUN echo "from c"
`

	assert.Equal(t, exp, out)
}
