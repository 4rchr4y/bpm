package main

import (
	"log"

	"github.com/4rchr4y/bpm/cli/require"
	"github.com/4rchr4y/bpm/fileifier"
	"github.com/4rchr4y/bpm/internal/encode"
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
	tomlEncoder := encode.NewTomlEncoder()

	bundleParser := fileifier.NewFileifier(tomlEncoder)
	gitService := gitcli.NewClient()
	gitLoader := loader.NewGitLoader(gitService, bundleParser)

	bpmClient.RegisterCommand(
		manager.NewInstallCommand(&manager.InstallCmdHub{
			OsWrap:          osWrap,
			BundleInstaller: manager.NewBundleInstaller(osWrap, tomlEncoder),
			FileLoader:      gitLoader,
		}),

		manager.NewGetCommand(&manager.GetCmdHub{
			OsWrap:      osWrap,
			TomlEncoder: tomlEncoder,
			FileLoader:  gitLoader,
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

	if _, err := manager.ExecuteGetCmd(getCmd, &manager.GetCmdInput{
		URL:     pathToBundle,
		Version: version,
	}); err != nil {
		log.Fatal(err)
		return
	}
}
