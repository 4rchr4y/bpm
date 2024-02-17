package get

import (
	"context"
	"os"

	"github.com/4rchr4y/bpm/bundle"
	"github.com/4rchr4y/bpm/bundleutil/encode"
	"github.com/4rchr4y/bpm/bundleutil/manifest"
	"github.com/4rchr4y/bpm/cli/cmdutil/factory"
	"github.com/4rchr4y/bpm/cli/cmdutil/require"
	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/fetch"
	"github.com/4rchr4y/bpm/storage"
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
				workDir:    wd,
				url:        args[0],
				Version:    version,
				Fetcher:    f.Fetcher,
				Storage:    f.Storage,
				Encoder:    f.Encoder,
				Manifester: f.Manifester,
			})
		},
	}

	cmd.Flags().StringP("version", "v", "", "Bundle version")
	return cmd
}

type getOptions struct {
	io         core.IO
	workDir    string // bundle working directory
	url        string // bundle repository that needs to be installed
	Version    string // specified bundle version
	Fetcher    *fetch.Fetcher
	Storage    *storage.Storage
	Encoder    *encode.Encoder      // decoder of bundle component files
	Manifester *manifest.Manifester // bundle manifest file control operator
}

func getRun(ctx context.Context, opts *getOptions) error {
	dest, err := opts.Storage.LoadFromAbs(opts.workDir, nil)
	if err != nil {
		return err
	}

	v, err := bundle.ParseVersionExpr(opts.Version)
	if err != nil {
		return err
	}

	input := &manifest.InsertRequirementInput{
		Parent:  dest,
		Source:  opts.url,
		Version: v,
	}

	if err := opts.Manifester.InsertRequirement(ctx, input); err != nil {
		return err
	}

	if err := opts.Manifester.Upgrade(opts.workDir, dest); err != nil {
		return err
	}

	return nil
}
