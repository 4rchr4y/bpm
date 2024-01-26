package version

import (
	"fmt"

	"github.com/4rchr4y/bpm/command/factory"
	"github.com/spf13/cobra"
)

func NewCmdVersion(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "version",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(f.Version)
		},
	}

	return cmd
}
