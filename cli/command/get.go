package command

import (
	"log"

	"github.com/4rchr4y/bpm"
	"github.com/4rchr4y/bpm/internal/domain/encode"
	gitcli "github.com/4rchr4y/bpm/internal/domain/service/git"
	"github.com/4rchr4y/godevkit/syswrap"
	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a new dependency",
	Long:  ``,
	// Args:  validateGetCmdArgs,
	Run: runGetCmd,
}

func init() {
	RootCmd.AddCommand(GetCmd)

	GetCmd.Flags().StringP("version", "v", "", "Bundle version")
}

func runGetCmd(cmd *cobra.Command, args []string) {
	pathToBundle := args[0]
	bpmClient := bpm.NewBpm()
	osWrap := new(syswrap.OsWrapper)
	tomlEncoder := encode.NewTomlEncoder()

	bundleParser := bpm.NewBundleParser(tomlEncoder)
	gitService := gitcli.NewClient()
	gitLoader := bpm.NewGitLoader(gitService, bundleParser)

	bpmClient.RegisterCommand(
		bpm.NewInstallCommand(&bpm.InstallCmdConf{
			OsWrap:      osWrap,
			TomlEncoder: tomlEncoder,
			FileLoader:  gitLoader,
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
