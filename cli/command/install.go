package command

import (
	"log"

	"github.com/4rchr4y/bpm"
	"github.com/4rchr4y/bpm/internal/domain/encode"
	gitcli "github.com/4rchr4y/bpm/internal/domain/service/git"
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
	bpmClient := bpm.NewBpm()
	osWrap := new(syswrap.OsWrapper)
	tomlEncoder := encode.NewTomlEncoder()

	bundleParser := bpm.NewBundleParser(tomlEncoder)
	gitService := gitcli.NewClient()
	gitLoader := bpm.NewGitLoader(gitService, bundleParser)

	bpmClient.RegisterCommand(
		bpm.NewInstallCommand(&bpm.InstallCmdConf{
			OsWrap:          osWrap,
			TomlEncoder:     tomlEncoder,
			BundleInstaller: bpm.NewBundleInstaller(osWrap, tomlEncoder),
			FileLoader:      gitLoader,
		}),
	)

	getCmd, err := bpmClient.Command(bpm.GetCommandName)
	if err != nil {
		log.Fatal(err)
		return
	}

	version, err := cmd.Flags().GetString("version")
	if err != nil {
		log.Fatal(err)
		return
	}

	if _, err := getCmd.Execute(&bpm.InstallCmdInput{
		URL:     pathToBundle,
		Version: version,
	}); err != nil {
		log.Fatal(err)
		return
	}
}
