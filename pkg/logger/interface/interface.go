package logger_interface

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
}
