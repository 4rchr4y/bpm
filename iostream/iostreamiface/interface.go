package iostreamiface

import (
	"io"

	"github.com/4rchr4y/bpm/core"
)

type IO interface {
	Println(a ...any)
	Printf(format string, a ...any)

	PrintfWarn(format string, a ...any)
	PrintfErr(format string, a ...any)
	PrintfDebug(format string, a ...any)
	PrintfOk(format string, a ...any)
	PrintfInfo(format string, a ...any)

	GetStdin() io.Reader
	GetStdout() io.Writer
	GetStdoutErr() io.Writer
	GetStdoutMode(mode core.StdoutMode) core.StdoutMode

	SetStdoutMode(mode core.StdoutMode)
}
