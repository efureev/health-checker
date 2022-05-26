package checker

type ILogger interface {
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Success(format string, args ...interface{})
	Error(format string, args ...interface{})
	Log(format string, args ...interface{})
}
