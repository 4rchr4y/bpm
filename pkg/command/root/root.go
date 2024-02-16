package root

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/pkg/cmdutil/factory"
	"github.com/spf13/cobra"

	cmdCheck "github.com/4rchr4y/bpm/pkg/command/check"
	cmdGet "github.com/4rchr4y/bpm/pkg/command/get"
	cmdInit "github.com/4rchr4y/bpm/pkg/command/init"
	cmdInstall "github.com/4rchr4y/bpm/pkg/command/install"
	cmdVersion "github.com/4rchr4y/bpm/pkg/command/version"
)

func NewCmdRoot(f *factory.Factory, version string) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:           "bpm",
		Short:         "Bundle Package Manager",
		Long:          "",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			debug, err := cmd.Flags().GetBool("debug")
			if err != nil {
				return err
			}

			if debug {
				f.IOStream.SetStdoutMode(core.Debug)
			}

			return nil
		},
	}

	cmd.AddGroup(&cobra.Group{
		ID:    "core",
		Title: "Core commands",
	})

	cmd.PersistentFlags().Bool("debug", false, "Run `bpm` in debug mode")

	cmd.AddCommand(cmdVersion.NewCmdVersion(f))
	cmd.AddCommand(cmdInit.NewCmdInit(f))
	cmd.AddCommand(cmdInstall.NewCmdInstall(f))
	cmd.AddCommand(cmdCheck.NewCmdCheck(f))
	cmd.AddCommand(cmdGet.NewCmdGet(f))

	return cmd, nil
}

func rootUsageFunc(w io.Writer, command *cobra.Command) error {
	// fmt.Fprintf(w, "Usage:  %s", command.UseLine())

	// var subcommands []*cobra.Command
	// for _, c := range command.Commands() {
	// 	if !c.IsAvailableCommand() {
	// 		continue
	// 	}
	// 	subcommands = append(subcommands, c)
	// }

	// if len(subcommands) > 0 {
	// 	fmt.Fprint(w, "\n\nAvailable commands:\n")
	// 	for _, c := range subcommands {
	// 		fmt.Fprintf(w, "  %s\n", c.Name())
	// 	}
	// 	return nil
	// }

	// flagUsages := command.LocalFlags().FlagUsages()
	// if flagUsages != "" {
	// 	fmt.Fprintln(w, "\n\nFlags:")
	// 	fmt.Fprint(w, text.Indent(dedent(flagUsages), "  "))
	// }
	return nil
}

func dedent(s string) string {
	lines := strings.Split(s, "\n")
	minIndent := -1

	for _, l := range lines {
		if len(l) == 0 {
			continue
		}

		indent := len(l) - len(strings.TrimLeft(l, " "))
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	if minIndent <= 0 {
		return s
	}

	var buf bytes.Buffer
	for _, l := range lines {
		fmt.Fprintln(&buf, strings.TrimPrefix(l, strings.Repeat(" ", minIndent)))
	}
	return strings.TrimSuffix(buf.String(), "\n")
}
