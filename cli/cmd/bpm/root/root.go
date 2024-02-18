package root

import (
	"github.com/4rchr4y/bpm/cli/cmdutil/factory"
	"github.com/4rchr4y/bpm/core"
	"github.com/spf13/cobra"

	cmdGet "github.com/4rchr4y/bpm/cli/cmd/bpm/get"
	cmdInit "github.com/4rchr4y/bpm/cli/cmd/bpm/init"
	cmdInstall "github.com/4rchr4y/bpm/cli/cmd/bpm/install"
	cmdTidy "github.com/4rchr4y/bpm/cli/cmd/bpm/tidy"
	cmdVersion "github.com/4rchr4y/bpm/cli/cmd/bpm/version"
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
	cmd.AddCommand(cmdTidy.NewCmdTidy(f))
	cmd.AddCommand(cmdGet.NewCmdGet(f))

	return cmd, nil
}
