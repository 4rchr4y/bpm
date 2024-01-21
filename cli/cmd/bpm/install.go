package main

import (
	"log"

	"github.com/4rchr4y/bpm/cli/require"
	"github.com/4rchr4y/bpm/fileifier"
	"github.com/4rchr4y/bpm/internal/encode"
	"github.com/4rchr4y/bpm/loader"
	"github.com/4rchr4y/bpm/manager"
	"github.com/4rchr4y/godevkit/syswrap"
	"github.com/spf13/cobra"

	gitcli "github.com/4rchr4y/bpm/internal/git"
)

func newInstallCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install [REPOSITORY]",
		Aliases: []string{"i"},
		Args:    require.ExactArgs(1),
		Short:   "Install a package from the specified repository",
		Run:     runInstallCmd,
	}

	cmd.Flags().StringP("version", "v", "", "Bundle version")
	return cmd
}

func runInstallCmd(cmd *cobra.Command, args []string) {
	pathToBundle := args[0]
	bpmClient := manager.NewBpm()
	osWrap := new(syswrap.OsWrapper)
	tomlEncoder := encode.NewTomlEncoder()

	bundleParser := fileifier.NewFileifier(tomlEncoder)
	gitService := gitcli.NewClient()
	gitLoader := loader.NewGitLoader(gitService, bundleParser)

	bpmClient.RegisterCommand(
		manager.NewInstallCommand(&manager.InstallCmdConf{
			OsWrap:          osWrap,
			BundleInstaller: manager.NewBundleInstaller(osWrap, tomlEncoder),
			FileLoader:      gitLoader,
		}),
	)

	getCmd, err := bpmClient.Command(manager.InstallCmdName)
	if err != nil {
		log.Fatal(err)
		return
	}

	version, err := cmd.Flags().GetString("version")
	if err != nil {
		log.Fatal(err)
		return
	}

	if _, err := getCmd.Execute(&manager.InstallCmdInput{
		URL:     pathToBundle,
		Version: version,
	}); err != nil {
		log.Fatal(err)
		return
	}
}
