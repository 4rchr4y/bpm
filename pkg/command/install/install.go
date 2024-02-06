package install

import (
	"context"

	"github.com/4rchr4y/bpm/pkg/bundleutil"
	"github.com/4rchr4y/bpm/pkg/cmdutil/factory"
	"github.com/4rchr4y/bpm/pkg/cmdutil/require"
	"github.com/spf13/cobra"
)

func NewCmdInstall(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install REPOSITORY",
		Aliases: []string{"i"},
		Args:    require.ExactArgs(1),
		Short:   "Install a package from the specified repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			version, err := cmd.Flags().GetString("version")
			if err != nil {
				return err
			}

			return installRun(cmd.Context(), &installOptions{
				URL:     args[0],
				Version: version,
				Loader:  f.Loader,
				Saver:   f.Saver,
			})
		},
	}

	cmd.Flags().StringP("version", "v", "", "Bundle version")
	return cmd
}

type installOptions struct {
	URL     string // bundle repository that needs to be installed
	Version string // specified bundle version
	Loader  *bundleutil.Loader
	Saver   *bundleutil.Saver // bundle saver files into the file system
}

func installRun(ctx context.Context, opts *installOptions) error {
	// b, err := opts.Loader.DownloadBundle(ctx, opts.URL, opts.Version)
	// if err != nil {
	// 	return err
	// }

	// if err := opts.Saver.SaveToDisk(b); err != nil {
	// 	return err
	// }

	return nil
}
