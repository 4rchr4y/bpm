package core

type IO interface {
	Println(a ...any)
	PrintfWarn(format string, a ...any)
	PrintfErr(format string, a ...any)
	Printf(format string, a ...any)
}
