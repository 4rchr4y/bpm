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
	bpmClient := manager.NewBpm()
	osWrap := new(syswrap.OsWrapper)
	tomlEncoder := encode.NewTomlEncoder()

	bundleParser := fileifier.NewFileifier(tomlEncoder)
	gitService := gitcli.NewClient()
	gitLoader := loader.NewGitLoader(gitService, bundleParser)

	bpmClient.RegisterCommand(
		manager.NewInstallCommand(&manager.InstallCmdConf{
			OsWrap:     osWrap,
			FileLoader: gitLoader,
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
