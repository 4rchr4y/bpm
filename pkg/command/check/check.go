package check

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/4rchr4y/bpm/cli/require"
	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/pkg/bundleutil"
	"github.com/4rchr4y/bpm/pkg/cmdutil/factory"
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
				Loader:    f.Loader,
			})
		},
	}

	return cmd
}

type checkOptions struct {
	Path      string                // path to the bundle that needs to be checked
	Encoder   *bundleutil.Encoder   // decoder of bundle component files
	Fileifier *bundleutil.Fileifier // transformer of file contents into structures
	Loader    *bundleutil.Loader    // bundle file loader from file system
}

func checkRun(opts *checkOptions) error {
	b, err := opts.Loader.LoadBundle(opts.Path)
	if err != nil {
		return fmt.Errorf("failed to load '%s' bundle: %v", opts.Path, err)
	}

	// for _, file := range b.RegoFiles {
	// 	b.LockFile.SetModule(&lockfile.ModDecl{
	// 		Package: file.Package(),
	// 		Source:  file.Path,
	// 		Sum:     file.Sum(),
	// 	})
	// }

	if err := os.WriteFile(
		filepath.Join(opts.Path, constant.LockFileName),
		opts.Encoder.EncodeLockFile(b.LockFile),
		0644,
	); err != nil {
		return fmt.Errorf("failed to write file '%s': %v", "fileName", err)
	}

	return nil
}
