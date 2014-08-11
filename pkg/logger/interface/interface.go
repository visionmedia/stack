package logger_interface

type Logger interface {
	Skip(arg, commit string)
	Start(arg, commit string)
	Success(arg, commit string)
	Error(err error, arg, commit string)
	Log(ling string)
	End()
}
