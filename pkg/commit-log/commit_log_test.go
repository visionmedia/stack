package commit_log

import "github.com/bmizerany/assert"
import "testing"
import "os"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Test new with no file.
func TestNewEmpty(t *testing.T) {
	os.Remove("/tmp/commits.log")

	c, err := New("/tmp/commits.log")
	check(err)

	assert.Equal(t, 0, c.Length())
	assert.Equal(t, false, c.Has("foo"))
}

// Test new with commits.
func TestNewCommits(t *testing.T) {
	os.Remove("/tmp/commits.log")

	c, err := New("/tmp/commits.log")
	check(err)

	check(c.Add("foo"))
	check(c.Add("bar"))
	check(c.Add("baz"))

	check(c.Close())
	check(c.Open())

	assert.Equal(t, 3, c.Length())
	assert.Equal(t, true, c.Has("foo"))
	assert.Equal(t, true, c.Has("bar"))
	assert.Equal(t, true, c.Has("baz"))
	assert.Equal(t, false, c.Has("something"))
}

// Test clearing the log.
func TestClear(t *testing.T) {
	os.Remove("/tmp/commits.log")

	c, err := New("/tmp/commits.log")
	check(err)

	check(c.Add("foo"))
	check(c.Add("bar"))
	check(c.Add("baz"))

	check(c.Clear())

	assert.Equal(t, 0, c.Length())
	assert.Equal(t, false, c.Has("foo"))
	assert.Equal(t, false, c.Has("bar"))
	assert.Equal(t, false, c.Has("baz"))
	assert.Equal(t, false, c.Has("something"))
}
