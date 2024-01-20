package command

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "bpm",
	Short: "Bundle Package Manager",
	Long:  "",
}
