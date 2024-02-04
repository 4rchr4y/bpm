package iostream

import (
	"fmt"
	"io"

	"github.com/4rchr4y/bpm/core"
	"github.com/muesli/termenv"
)

var (
	LabelErr  = termenv.String(" ERROR ").Bold().Background(DarkThemeRedDeep).String()
	LabelWarn = termenv.String(" WARN ").Bold().Background(DarkThemeYellowLight).String()
)

type IOStream struct {
	mode core.PrintMode
	in   io.Reader
	out,
	errOut *termenv.Output
}

func NewIOStream(in io.Reader, out io.Writer, outErr io.Writer, mode ...core.PrintMode) *IOStream {
	if len(mode) < 1 {
		mode[0] = core.Info
	}

	return &IOStream{
		mode:   mode[0],
		in:     in,
		out:    termenv.NewOutput(out),
		errOut: termenv.NewOutput(outErr),
	}
}

func (s *IOStream) SetPrintMode(mode core.PrintMode) { s.mode = mode }

func (s *IOStream) Println(a ...any)               { fmt.Fprintln(s.out, a...) }
func (s *IOStream) Printf(format string, a ...any) { fmt.Fprintf(s.out, format, a...) }

func (s *IOStream) PrintfWarn(format string, a ...any) {
	msg := termenv.String(fmt.Sprintf(format, a...)).Foreground(DarkThemeYellowLight).String()
	str := termenv.String(LabelWarn, msg).String()

	fmt.Fprint(s.out, str+"\n")
}

func (s *IOStream) PrintfErr(format string, a ...any) {
	msg := termenv.String(fmt.Sprintf(format, a...)).Foreground(DarkThemeRedDeep).String()
	str := termenv.String(LabelErr, msg).String()

	fmt.Fprint(s.errOut, str+"\n")
}
