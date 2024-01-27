package root

import (
	"github.com/4rchr4y/bpm/pkg/command/factory"
	"github.com/spf13/cobra"

	cmdCheck "github.com/4rchr4y/bpm/pkg/command/check"
	cmdGet "github.com/4rchr4y/bpm/pkg/command/get"
	cmdInit "github.com/4rchr4y/bpm/pkg/command/init"
	cmdInstall "github.com/4rchr4y/bpm/pkg/command/install"
	cmdVersion "github.com/4rchr4y/bpm/pkg/command/version"
)

func NewCmdRoot(f *factory.Factory, version string) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "bpm",
		Short: "Bundle Package Manager",
		Long:  "",
	}

	cmd.AddGroup(&cobra.Group{
		ID:    "core",
		Title: "Core commands",
	})

	cmd.AddCommand(cmdVersion.NewCmdVersion(f))
	cmd.AddCommand(cmdInit.NewCmdInit(f))
	cmd.AddCommand(cmdInstall.NewCmdInstall(f))
	cmd.AddCommand(cmdCheck.NewCmdCheck(f))
	cmd.AddCommand(cmdGet.NewCmdGet(f))

	return cmd, nil
}
