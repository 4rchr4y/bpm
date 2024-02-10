package check

import (
	"github.com/4rchr4y/bpm/pkg/bundleutil/encode"
	"github.com/4rchr4y/bpm/pkg/cmdutil/factory"
	"github.com/4rchr4y/bpm/pkg/cmdutil/require"
	"github.com/4rchr4y/bpm/pkg/fetch"
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
				Path:    args[0],
				Encoder: f.Encoder,
				Fetcher: f.Fetcher,
			})
		},
	}

	return cmd
}

type checkOptions struct {
	Path    string          // path to the bundle that needs to be checked
	Encoder *encode.Encoder // decoder of bundle component files
	Fetcher *fetch.Fetcher  // bundle file loader from file system
}

func checkRun(opts *checkOptions) error {
	// b, err := opts.Fetcher.FetchLocal(opts.Path)
	// if err != nil {
	// 	return fmt.Errorf("failed to load '%s' bundle: %v", opts.Path, err)
	// }

	// if err := os.WriteFile(
	// 	filepath.Join(opts.Path, constant.LockFileName),
	// 	opts.Encoder.EncodeLockFile(b.LockFile),
	// 	0644,
	// ); err != nil {
	// 	return fmt.Errorf("failed to write file '%s': %v", "fileName", err)
	// }

	return nil
}
