package main

import (
	"github.com/4rchr4y/bpm/cli/require"
	"github.com/4rchr4y/bpm/fileifier"
	"github.com/4rchr4y/bpm/internal/encode"
	"github.com/4rchr4y/bpm/loader"
	"github.com/4rchr4y/bpm/manager"
	"github.com/4rchr4y/godevkit/syswrap"
	"github.com/spf13/cobra"
)

func newCheckCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "check PATH",
		Aliases: []string{"validate"},
		Args:    require.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				// Allow file completion when completing the argument for the name
				// which could be a path
				return nil, cobra.ShellCompDirectiveDefault
			}

			// No more completions, so disable file completion
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Short: "Check the validity of the bundle using the specified path",
		RunE:  runCheckCmd,
	}

	return cmd
}

func runCheckCmd(cmd *cobra.Command, args []string) error {
	pathToBundle := args[0]
	bpmManager := manager.NewBpm()
	osWrap := new(syswrap.OsWrapper)
	ioWrap := new(syswrap.IoWrapper)
	tomlEncoder := encode.NewTomlEncoder()

	fileifier := fileifier.NewFileifier(tomlEncoder)
	fileLoader := loader.NewFsLoader(osWrap, ioWrap, fileifier)

	checkCmd := manager.NewCheckCommand(&manager.CheckCmdResources{
		FileLoader: fileLoader,
	})

	if err := bpmManager.RegisterCommand(
		checkCmd,
	); err != nil {
		return err
	}

	if _, err := manager.ExecuteCheckCmd(checkCmd, &manager.CheckCmdInput{
		Path: pathToBundle,
	}); err != nil {
		return err
	}

	return nil
}
