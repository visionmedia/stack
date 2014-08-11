//
// Commit log which stores sha1 hashes on disk
// so that re-runs of the same command may be
// ignored.
//
package commit_log

import "path/filepath"
import "strings"
import "bufio"
import "sync"
import "io"
import "os"

type CommitLog struct {
	Path    string
	commits map[string]bool
	file    *os.File
	sync.Mutex
}

// New commit log stored at the given `path`.
func New(path string) (*CommitLog, error) {
	c := &CommitLog{
		Path:    path,
		commits: make(map[string]bool),
	}

	err := c.Open()
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Open the commit log and read its contents.
// The directory will be created if it does not exist.
func (c *CommitLog) Open() error {
	dir := filepath.Dir(c.Path)

	err := os.MkdirAll(dir, 0644)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(c.Path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	c.file = file
	r := bufio.NewReader(file)

	for {
		commit, err := r.ReadString('\n')

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		commit = strings.Trim(commit, "\n")
		c.commits[commit] = true
	}

	return nil
}

// Clear the commit log and re-open.
func (c *CommitLog) Clear() error {
	err := c.Close()
	if err != nil {
		return err
	}

	err = os.Remove(c.Path)
	if err != nil {
		return err
	}

	return c.Open()
}

// Close the commit log.
func (c *CommitLog) Close() error {
	c.Lock()
	defer c.Unlock()
	c.commits = make(map[string]bool)
	return c.file.Close()
}

// Add `commit` and write to disk.
func (c *CommitLog) Add(commit string) error {
	c.Lock()
	defer c.Unlock()

	c.commits[commit] = true

	_, err := io.WriteString(c.file, commit+"\n")
	return err
}

// Has checks if `commit` exists.
func (c *CommitLog) Has(commit string) bool {
	c.Lock()
	defer c.Unlock()
	return c.commits[commit]
}

// Length of commit log.
func (c *CommitLog) Length() int {
	c.Lock()
	defer c.Unlock()
	return len(c.commits)
}
