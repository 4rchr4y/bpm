package command

import (
	"log"

	"github.com/4rchr4y/bpm/fileifier"
	"github.com/4rchr4y/bpm/internal/encode"
	gitcli "github.com/4rchr4y/bpm/internal/git"
	"github.com/4rchr4y/bpm/loader"
	"github.com/4rchr4y/bpm/manager"
	"github.com/4rchr4y/godevkit/syswrap"
	"github.com/spf13/cobra"
)

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a new dependency",
	Long:  ``,
	Run:   runInstallCmd,
}

func init() {
	RootCmd.AddCommand(InstallCmd)

	InstallCmd.Flags().StringP("version", "v", "", "Bundle version")
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

	getCmd, err := bpmClient.Command(manager.GetCommandName)
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
