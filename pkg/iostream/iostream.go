package iostream

import (
	"fmt"
	"io"

	"github.com/muesli/termenv"
)

var (
	LabelErr  = termenv.String(" ERROR ").Bold().Background(DarkThemeRedDeep).String()
	LabelWarn = termenv.String(" WARN ").Bold().Background(DarkThemeYellowLight).String()
)

type IOStream struct {
	in          io.Reader
	Out, ErrOut *termenv.Output
}

func NewIOStream(in io.Reader, out io.Writer, outErr io.Writer) *IOStream {
	return &IOStream{
		in:     in,
		Out:    termenv.NewOutput(out),
		ErrOut: termenv.NewOutput(outErr),
	}
}

func (s *IOStream) Println(a ...any)               { fmt.Fprintln(s.Out, a...) }
func (s *IOStream) Printf(format string, a ...any) { fmt.Fprintf(s.Out, format, a...) }

func (s *IOStream) PrintfWarn(format string, a ...any) {
	msg := termenv.String(fmt.Sprintf(format, a...)).Foreground(DarkThemeYellowLight).String()
	str := termenv.String(LabelWarn, msg).String()

	fmt.Fprint(s.Out, str+"\n")
}

func (s *IOStream) PrintfErr(format string, a ...any) {
	msg := termenv.String(fmt.Sprintf(format, a...)).Foreground(DarkThemeRedDeep).String()
	str := termenv.String(LabelErr, msg).String()

	fmt.Fprint(s.ErrOut, str+"\n")
}
