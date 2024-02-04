package iostream

import (
	"fmt"
	"io"
	"os"

	"github.com/4rchr4y/bpm/core"
	"github.com/muesli/termenv"
)

var (
	LabelErr   = termenv.String(" ERROR ").Bold().Background(DarkThemeRedDeep).String()
	LabelDebug = termenv.String(" DEBUG ").Bold().Background(DarkThemeOrangeLight).String()
	LabelWarn  = termenv.String(" WARN ").Bold().Background(DarkThemeYellowLight).String()
	LabelOk    = termenv.String(" OK ").Bold().Background(DarkThemeGreen).String()
)

type IOStream struct {
	in io.Reader
	out,
	errOut *termenv.Output
	mode core.StdoutMode
}

type IOStreamOptFn func(io *IOStream)

func WithInput(input io.Reader) IOStreamOptFn {
	return func(io *IOStream) {
		io.in = input
	}
}

func WithOutput(output io.Writer) IOStreamOptFn {
	return func(io *IOStream) {
		io.out = termenv.NewOutput(output)
	}
}

func WithErrOutput(errOutput io.Writer) IOStreamOptFn {
	return func(io *IOStream) {
		io.errOut = termenv.NewOutput(errOutput)
	}
}

func WithMode(mode core.StdoutMode) IOStreamOptFn {
	return func(io *IOStream) {
		io.mode = mode
	}
}

func NewIOStream(options ...IOStreamOptFn) *IOStream {
	io := &IOStream{
		mode:   core.Info,
		in:     os.Stdin,
		out:    termenv.NewOutput(os.Stdout),
		errOut: termenv.NewOutput(os.Stderr),
	}

	for _, optionFn := range options {
		optionFn(io)
	}

	return io
}

func (s *IOStream) GetStdoutMode(mode core.StdoutMode) core.StdoutMode { return s.mode }
func (s *IOStream) SetStdoutMode(mode core.StdoutMode)                 { s.mode = mode }

func (s *IOStream) Println(a ...any)               { fmt.Fprintln(s.out, a...) }
func (s *IOStream) Printf(format string, a ...any) { fmt.Fprintf(s.out, format, a...) }

func (s *IOStream) PrintfOk(format string, a ...any) {
	msg := termenv.String(fmt.Sprintf(format, a...)).String()
	str := termenv.String(LabelOk, msg).String()

	fmt.Fprint(s.out, str+"\n")
}

func (s *IOStream) PrintfWarn(format string, a ...any) {
	msg := termenv.String(fmt.Sprintf(format, a...)).String()
	str := termenv.String(LabelWarn, msg).String()

	fmt.Fprint(s.out, str+"\n")
}

func (s *IOStream) PrintfDebug(format string, a ...any) {
	if s.mode != core.Debug {
		return
	}

	msg := termenv.String(fmt.Sprintf(format, a...)).String()
	str := termenv.String(LabelDebug, msg).String()

	fmt.Fprint(s.out, str+"\n")
}

func (s *IOStream) PrintfErr(format string, a ...any) {
	msg := termenv.String(fmt.Sprintf(format, a...)).Foreground(DarkThemeRedDeep).String()
	str := termenv.String(LabelErr, msg).String()

	fmt.Fprint(s.errOut, str+"\n")
}
