package logger_interface

import "io"

type Logger interface {
	// Skip a commit.
	Skip(arg, commit string)

	// Start a commit.
	Start(arg, commit string)

	// Successful commit.
	Success(arg, commit string)

	// Error executing a commit.
	Error(err error, arg, commit string)

	// Log line (comment).
	Log(line string)

	// End of provision.
	End()

	// Stdout writer.
	Stdout() io.Writer

	// Stderr writer.
	Stderr() io.Writer
}
