package main

import (
	"io"

	"github.com/spf13/cobra"
)

func newRootCmd(out io.Writer, args []string) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "bpm",
		Short: "Bundle Package Manager",
		Long:  "",
	}

	cmd.AddCommand(
		newInstallCmd(args),
		newGetCmd(args),
		newInitCmd(args),
		newCheckCmd(args),
	)

	return cmd, nil
}
