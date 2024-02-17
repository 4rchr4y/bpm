package download

import (
	"context"

	"github.com/4rchr4y/bpm/bundleutil/inspect"
	"github.com/4rchr4y/bpm/cli/cmdutil/factory"
	"github.com/4rchr4y/bpm/cli/cmdutil/require"
	"github.com/4rchr4y/bpm/fetch"
	"github.com/4rchr4y/bpm/storage"
	"github.com/spf13/cobra"
)

func NewCmdDownload(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "download",
		Aliases: []string{"d"},
		Args:    require.NoArgs,
		Short:   "Download all bundle requirements",
		RunE: func(cmd *cobra.Command, args []string) error {
			return downloadRun(cmd.Context(), &downloadOptions{
				storage: f.Storage,
				fetcher: f.Fetcher,
			})
		},
	}

	return cmd
}

type downloadOptions struct {
	storage   *storage.Storage
	fetcher   *fetch.Fetcher
	inspector *inspect.Inspector
}

func downloadRun(ctx context.Context, opts *downloadOptions) error {
	b, err := opts.storage.LoadFromAbs(".", nil)
	if err != nil {
		return err
	}

	if err := opts.inspector.Inspect(b); err != nil {
		return err
	}

	// if err := opts.Saver.SaveToDisk(b); err != nil {
	// 	return err
	// }

	return nil
}
