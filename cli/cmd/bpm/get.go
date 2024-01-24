package main

import (
	"log"
	"os"

	"github.com/4rchr4y/bpm/bfencoder"
	"github.com/4rchr4y/bpm/cli/require"
	"github.com/4rchr4y/bpm/fileifier"
	gitcli "github.com/4rchr4y/bpm/internal/git"
	"github.com/4rchr4y/bpm/loader"
	"github.com/4rchr4y/bpm/manager"
	"github.com/4rchr4y/godevkit/syswrap"
	"github.com/spf13/cobra"
)

func newGetCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a new dependency",
		Args:  require.ExactArgs(1),
		Run:   runGetCmd,
	}

	cmd.Flags().StringP("version", "v", "", "Bundle version")
	return cmd
}

func runGetCmd(cmd *cobra.Command, args []string) {
	pathToBundle := args[0]
	bpmClient := manager.NewBpm()
	osWrap := new(syswrap.OsWrapper)
	ioWrap := new(syswrap.IoWrapper)
	// tomlEncoder := encode.NewTomlEncoder()

	bfEncoder := bfencoder.NewEncoder()

	filef := fileifier.NewFileifier(bfEncoder)
	gitLoader := loader.NewGitLoader(gitcli.NewClient(), filef)
	fsLoader := loader.NewFsLoader(osWrap, ioWrap, filef)

	bpmClient.RegisterCommand(
		manager.NewInstallCommand(&manager.InstallCmdResources{
			OsWrap:          osWrap,
			BundleInstaller: manager.NewBundleInstaller(osWrap, bfEncoder),
			FileLoader:      gitLoader,
		}),

		manager.NewGetCommand(&manager.GetCmdResources{
			OsWrap:    osWrap,
			Encoder:   bfEncoder,
			GitLoader: gitLoader,
			FsLoader:  fsLoader,
		}),
	)

	getCmd, err := bpmClient.Command(manager.GetCmdName)
	if err != nil {
		log.Fatal(err)
		return
	}

	version, err := cmd.Flags().GetString("version")
	if err != nil {
		log.Fatal(err)
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}

	if _, err := manager.ExecuteGetCmd(getCmd, &manager.GetCmdInput{
		URL:     pathToBundle,
		Version: version,
		Dir:     wd,
	}); err != nil {
		log.Fatal(err)
		return
	}
}
