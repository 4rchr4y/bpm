package iostream

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/4rchr4y/bpm/core"
	"github.com/muesli/termenv"
)

var (
	LabelErr   = labelTemplate(termenv.String(" ERROR ").Foreground(DarkThemeRedDeep).String())
	LabelDebug = labelTemplate(termenv.String(" DEBUG ").Foreground(DarkThemeOrangeLight).String())
	LabelWarn  = labelTemplate(termenv.String(" WARN ").Foreground(DarkThemeYellowLight).String())
	LabelInfo  = labelTemplate(termenv.String(" INFO ").Foreground(DarkThemeFg0).String())
	LabelOk    = labelTemplate(termenv.String(" OK ").Foreground(DarkThemeGreen).String())
)

type IOStream struct {
	in io.Reader
	out,
	errOut *termenv.Output
	mode core.StdoutMode
	poet time.Time // execution time of the previous operation
}

type IOStreamOptFn func(io *IOStream)

func WithInput(input io.Reader) IOStreamOptFn {
	return func(io *IOStream) { io.in = input }
}

func WithOutput(output io.Writer) IOStreamOptFn {
	return func(io *IOStream) { io.out = termenv.NewOutput(output) }
}

func WithErrOutput(errOutput io.Writer) IOStreamOptFn {
	return func(io *IOStream) { io.errOut = termenv.NewOutput(errOutput) }
}

func WithMode(mode core.StdoutMode) IOStreamOptFn {
	return func(io *IOStream) { io.mode = mode }
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
	s.Println(s.prepareWithLabel(LabelOk, format, a...))
}

func (s *IOStream) PrintfInfo(format string, a ...any) {
	s.Println(s.prepareWithLabel(LabelInfo, format, a...))
}

func (s *IOStream) PrintfWarn(format string, a ...any) {
	s.Println(s.prepareWithLabel(LabelWarn, format, a...))
}

func (s *IOStream) PrintfDebug(format string, a ...any) {
	if s.mode != core.Debug {
		return
	}

	s.Println(s.prepareWithLabel(LabelDebug, format, a...))
}

func (s *IOStream) PrintfErr(format string, a ...any) {
	msg := termenv.String(fmt.Sprintf(format, a...)).Foreground(DarkThemeRedDeep).String()
	str := termenv.String(LabelErr, msg).String()

	fmt.Fprint(s.errOut, str+"\n")
}

func (s *IOStream) fetchTime(now time.Time) (result string) {
	result = termenv.String(now.Format(core.StdoutTimeFormat)).Faint().String()
	if s.poet.IsZero() || s.mode == core.Info {
		s.poet = now
		return result
	}

	s.poet = now
	duration := time.Since(s.poet)
	result += termenv.String(fmt.Sprintf(" (%s)", duration)).Faint().String()

	return result
}

func (s *IOStream) prepareWithLabel(label, format string, a ...any) string {
	msg := termenv.String(fmt.Sprintf(format, a...)).String()
	return termenv.String(label, msg, s.fetchTime(time.Now())).String()
}

func labelTemplate(labelText string) string {
	return termenv.String("[").Foreground(DarkThemeFg0).Bold().String() + labelText + termenv.String("]").Foreground(DarkThemeFg0).Bold().String()
}
