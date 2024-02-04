package core

type PrintMode int

const (
	Debug = iota << 1
	Info
)

type IO interface {
	Println(a ...any)
	PrintfWarn(format string, a ...any)
	PrintfErr(format string, a ...any)
	Printf(format string, a ...any)
}
