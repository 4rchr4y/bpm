package core

type StdoutMode int

const (
	Debug = iota << 1
	Info
)

const StdoutTimeFormat = "15:04:05"

type IO interface {
	Println(a ...any)
	Printf(format string, a ...any)

	PrintfWarn(format string, a ...any)
	PrintfErr(format string, a ...any)
	PrintfDebug(format string, a ...any)
	PrintfOk(format string, a ...any)
	PrintfInfo(format string, a ...any)

	GetStdoutMode(mode StdoutMode) StdoutMode
	SetStdoutMode(mode StdoutMode)
}
