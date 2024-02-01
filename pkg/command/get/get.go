package get

import (
	"context"
	"os"

	"github.com/4rchr4y/bpm/cli/require"
	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/bpm/pkg/bundleutil"
	"github.com/4rchr4y/bpm/pkg/cmdutil/factory"
	"github.com/4rchr4y/bpm/pkg/load/osload"
	"github.com/spf13/cobra"
)

func NewCmdGet(f *factory.Factory) *cobra.Command {
	ctx := context.Background()
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a new dependency",
		Args:  require.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			version, err := cmd.Flags().GetString("version")
			if err != nil {
				return err
			}

			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			return getRun(ctx, &getOptions{
				WorkDir:    wd,
				URL:        args[0],
				Version:    version,
				OsLoader:   f.OsLoader,
				Saver:      f.Saver,
				Encoder:    f.Encoder,
				Downloader: f.Downloader,
				Manifester: f.Manifester,
			})
		},
	}

	cmd.Flags().StringP("version", "v", "", "Bundle version")
	return cmd
}

type getOptions struct {
	WorkDir    string                 // bundle working directory
	URL        string                 // bundle repository that needs to be installed
	Version    string                 // specified bundle version
	OsLoader   *osload.OsLoader       // bundle file loader from file system
	Saver      *bundleutil.Saver      // bundle saver files into the file system
	Encoder    *bundleutil.Encoder    // decoder of bundle component files
	Downloader *bundleutil.Downloader // downloader of a bundle and its dependencies
	Manifester *bundleutil.Manifester // bundle manifest file control operator
}

func getRun(ctx context.Context, opts *getOptions) error {
	target, err := opts.OsLoader.LoadBundle(opts.WorkDir)
	if err != nil {
		return err
	}

	result, err := opts.Downloader.Download(ctx, opts.URL, opts.Version)
	if err != nil {
		return err
	}

	updateInput := &bundleutil.UpdateInput{
		Target:    target,
		Rdirect:   append([]*bundle.Bundle(nil), result.Target),
		Rindirect: append(result.Rdirect, result.Rindirect...),
	}

	if err := opts.Manifester.Update(updateInput); err != nil {
		return err
	}

	if err := opts.Manifester.Upgrade(opts.WorkDir, target); err != nil {
		return err
	}

	return nil
}
