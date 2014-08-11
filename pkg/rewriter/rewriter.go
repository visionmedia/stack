//
// Rewriter accepts an input reader and pre-processes
// the shell-like provisioning language to its
// canonical form. For example command lines become
// prefixed with "RUN ", and comments prefixed with "LOG "
// and so on.
//
package rewriter

import "strings"
import "bufio"
import "fmt"
import "io"
import "os"

// Rewrite the given reader.
func Rewrite(r io.Reader) (string, error) {
	buf := bufio.NewReader(r)
	ret := ""

	for {
		line, readError := buf.ReadString('\n')

		str, err := rewrite(line)
		if err != nil {
			return "", err
		}

		ret += str

		if readError == io.EOF {
			break
		}

		if readError != nil {
			return "", readError
		}
	}

	return ret, nil
}

// Rewrite the given line.
func rewrite(line string) (string, error) {
	line = strings.Trim(line, " \r\n")

	switch {
	case "" == line:
		return "", nil
	case isLog(line):
		return log(line), nil
	case isInclude(line):
		return include(line)
	case isRun(line):
		return run(line), nil
	}

	return "", nil
}

// RUN line.
func run(line string) string {
	if hasCommandPrefix(line, "run") {
		return "RUN " + strip(line) + "\n"
	} else {
		return "RUN " + line + "\n"
	}
}

// LOG line.
func log(line string) string {
	return "LOG " + strip(line) + "\n"
}

// Strip the command.
func strip(line string) string {
	parts := strings.SplitN(line, " ", 2)
	return strings.Trim(parts[1], " ")
}

// Check if it's a RUN command.
func isRun(line string) bool {
	return hasCommandPrefix(line, "run") || true
}

// Check if it's an INCLUDE command.
func isInclude(line string) bool {
	return hasCommandPrefix(line, ".") || hasCommandPrefix(line, "include") || hasCommandPrefix(line, "source")
}

// Check if it's a LOG command.
func isLog(line string) bool {
	return (len(line) > 0 && line[0] == '#') || hasCommandPrefix(line, "log")
}

// Check for a command prefix.
func hasCommandPrefix(line, cmd string) bool {
	return strings.HasPrefix(strings.ToLower(line), cmd+" ")
}

// Include the file at `path` and rewrite it.
func include(path string) (string, error) {
	path = strip(path)

	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to include %q: %s", path, err)
	}

	return Rewrite(f)
}
