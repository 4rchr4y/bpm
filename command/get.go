package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Install a new dependency",
	Long:  ``,
	// Args:  validateGetCmdArgs,
	Run: runGetCmd,
}

func runGetCmd(cmd *cobra.Command, args []string) {
	fmt.Println("Hello, World!")
}
