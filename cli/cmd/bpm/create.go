package main

import (
	"github.com/4rchr4y/bpm/cli/require"
	"github.com/4rchr4y/godevkit/syswrap"
	"github.com/spf13/cobra"
)

const newCmdDesc = `
The 'bpm new' command is designed to generate a new bundle directory 
structure, complete with standard files typical for a bundle.

For instance, executing bpm new foo would result in a directory 
setup similar to the following:

    foo/
    ├── .bpmignore   	# File with patterns to exclude in bundle packaging.
    ├── bundle.toml    	# File with bundle information.
    ├── bpm.work   		# File with common information.
    └── .gitignore      # Ignore for git system.

'bpm new' takes a path for an argument. If directories in the given path
do not exist, bpm will attempt to create them as it goes. If the given
destination exists and there are files in that directory, conflicting files
will be overwritten, but other files will be left alone.
`

func newNewCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "new NAME",
		Aliases: []string{"n", "create"},
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
		Short: "Create a new bundle with the given name",
		Long:  newCmdDesc,
		RunE:  runNewCmd,
	}

	cmd.Flags().StringP("author", "a", "", "bundle author")
	return cmd
}

func runNewCmd(cmd *cobra.Command, args []string) error {
	osWrap := new(syswrap.OsWrapper)

	if err := osWrap.MkdirAll(args[0], 0755); err != nil {
		return err
	}

	return nil
}
