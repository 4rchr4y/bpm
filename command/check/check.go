package check

import (
	"fmt"

	"github.com/4rchr4y/bpm/cli/require"
	"github.com/4rchr4y/bpm/command/factory"
	"github.com/4rchr4y/bpm/fileifier"
	"github.com/4rchr4y/bpm/pkg/encode"
	"github.com/4rchr4y/bpm/pkg/load/osload"
	"github.com/spf13/cobra"
)

func NewCmdCheck(f *factory.Factory) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return checkRun(&checkOptions{
				Path:      args[0],
				Encoder:   f.Encoder,
				Fileifier: f.Fileifier,
				OsLoader:  f.OsLoader,
			})
		},
	}

	return cmd
}

type checkOptions struct {
	Path      string                // path to the bundle that needs to be checked
	Encoder   *encode.BundleEncoder // decoder of bundle component files
	Fileifier *fileifier.Fileifier  // transformer of file contents into structures
	OsLoader  *osload.OsLoader      // bundle file loader from file system
}

func checkRun(opts *checkOptions) error {
	b, err := opts.OsLoader.LoadBundle(opts.Path)
	if err != nil {
		return fmt.Errorf("failed to load '%s' bundle: %v", opts.Path, err)
	}

	fmt.Println(b.Name(), "OK")
	return nil
}
