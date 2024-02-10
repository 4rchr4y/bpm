package get

import (
	"context"
	"os"

	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/bpm/pkg/bundleutil/encode"
	"github.com/4rchr4y/bpm/pkg/bundleutil/manifest"
	"github.com/4rchr4y/bpm/pkg/cmdutil/factory"
	"github.com/4rchr4y/bpm/pkg/cmdutil/require"
	"github.com/4rchr4y/bpm/pkg/download"
	"github.com/4rchr4y/bpm/pkg/fetch"
	"github.com/4rchr4y/bpm/pkg/storage"
	"github.com/spf13/cobra"
)

func NewCmdGet(f *factory.Factory) *cobra.Command {
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

			return getRun(cmd.Context(), &getOptions{
				io:         f.IOStream,
				WorkDir:    wd,
				URL:        args[0],
				Version:    version,
				Fetcher:    f.Fetcher,
				Storage:    f.Storage,
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
	io         core.IO
	WorkDir    string // bundle working directory
	URL        string // bundle repository that needs to be installed
	Version    string // specified bundle version
	Fetcher    *fetch.Fetcher
	Storage    *storage.Storage
	Encoder    *encode.Encoder      // decoder of bundle component files
	Downloader *download.Downloader // downloader of a bundle and its dependencies
	Manifester *manifest.Manifester // bundle manifest file control operator
}

func getRun(ctx context.Context, opts *getOptions) error {
	target, err := opts.Storage.LoadFromAbs(opts.WorkDir)
	if err != nil {
		return err
	}

	v, err := bundle.ParseVersionExpr(opts.Version)
	if err != nil {
		return err
	}

	if target.BundleFile.IsRequirementListed(opts.URL, v.String()) {
		opts.io.PrintfOk("bundle '%s+%s' is already installed", opts.URL, v.String())
		return nil
	}

	result, err := opts.Downloader.DownloadWithContext(ctx, opts.URL, v)
	if err != nil {
		return err
	}

	if err := opts.Storage.StoreMultiple(result.Merge()); err != nil {
		return err
	}

	updateInput := &manifest.UpdateInput{
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

	opts.io.PrintfOk("bundle '%s+%s' was successfully added", result.Target.Repository(), result.Target.Version)

	return nil
}
