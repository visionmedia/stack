//
// Provisioner library which accepts a shell-like language
// or its canonical form for simple but effective provisioning.
//
package provisioner

import "github.com/visionmedia/stack/pkg/logger/interface"
import "github.com/visionmedia/stack/pkg/commit-log"
import "github.com/visionmedia/stack/pkg/rewriter"
import "encoding/hex"
import "crypto/sha1"
import "strings"
import "os/exec"
import "bytes"
import "bufio"
import "fmt"
import "os"
import "io"

const Version = "0.0.1"

type Provisioner struct {
	Log     logger_interface.Logger
	DryRun  bool
	Verbose bool
	path    string
	commits *commit_log.CommitLog
	r       io.Reader
}

// New provisioner from the given reader.
func New(r io.Reader, path string) *Provisioner {
	return &Provisioner{path: path, r: r}
}

// Run pending commands.
//
// Each command is parsed and checked against
// the set of committed or previously executed
// commands via sha1.
//
func (p *Provisioner) Run() error {
	in, err := rewriter.Rewrite(p.r)
	if err != nil {
		return err
	}

	b := bytes.NewBuffer([]byte(in))
	buf := bufio.NewReader(b)

	commits, err := commit_log.New(p.path)
	if err != nil {
		return err
	}

	p.commits = commits

	for {
		line, err := buf.ReadString('\n')

		if err == io.EOF {
			err = p.process(line)
			if err != nil {
				return err
			}
			break
		}

		err = p.process(line)
		if err != nil {
			return err
		}
	}

	p.Log.End()

	return nil
}

// process the given `line`.
func (p *Provisioner) process(line string) error {
	line = strings.Trim(line, " \r\n")

	if line == "" {
		return nil
	}

	cmd, arg, err := parse(line)
	if err != nil {
		return err
	}

	hash := sha1.Sum([]byte(line))
	commit := hex.EncodeToString(hash[:])

	return p.Command(cmd, arg, commit)
}

// Command runs the given `arg` against `cmd`.
func (p *Provisioner) Command(cmd, arg, commit string) error {
	switch strings.ToLower(cmd) {
	case "run":
		return p.CommandRun(arg, commit)
	case "log":
		return p.CommandLog(arg, commit)
	default:
		return fmt.Errorf("invalid command %q", cmd)
	}
}

// CommandRun implements the "RUN" command.
func (p *Provisioner) CommandRun(arg, commit string) error {
	if p.commits.Has(commit) {
		p.Log.Skip(arg, commit)
		return nil
	}

	p.Log.Start(arg, commit)

	if p.DryRun {
		p.Log.Success(arg, commit)
		return nil
	}

	c := exec.Command("sh", "-c", arg)
	c.Stdin = os.Stdin

	if p.Verbose {
		c.Stdout = p.Log.Stdout()
		c.Stderr = p.Log.Stderr()
	}

	err := c.Run()

	if err == nil {
		p.Log.Success(arg, commit)
		return p.commits.Add(commit)
	}

	p.Log.Error(err, arg, commit)
	return nil
}

// CommandLog implements the "LOG" command.
func (p *Provisioner) CommandLog(arg, commit string) error {
	p.Log.Log(arg)
	return nil
}

// parse and arguments from `line`.
func parse(line string) (string, string, error) {
	parts := strings.SplitN(line, " ", 2)

	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid command %q", line)
	}

	return parts[0], parts[1], nil
}
