package get

import (
	"context"
	"os"

	"github.com/4rchr4y/bpm/bundle"
	"github.com/4rchr4y/bpm/bundleutil/manifest"
	"github.com/4rchr4y/bpm/cli/cmdutil/factory"
	"github.com/4rchr4y/bpm/cli/cmdutil/require"
	"github.com/4rchr4y/bpm/iostream/iostreamiface"
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
				version:    version,
				storage:    f.Storage,
				manifester: f.Manifester,
			})
		},
	}

	cmd.Flags().StringP("version", "v", "", "Bundle version")
	return cmd
}

type getOptions struct {
	io         iostreamiface.IO
	workDir    string // bundle working directory
	url        string // bundle repository that needs to be installed
	version    string // specified bundle version
	storage    *storage.Storage
	manifester *manifest.Manifester // bundle manifest file control operator
}

func getRun(ctx context.Context, opts *getOptions) error {
	dest, err := opts.storage.LoadFromAbs(opts.workDir, nil)
	if err != nil {
		return err
	}

	v, err := bundle.ParseVersionExpr(opts.version)
	if err != nil {
		return err
	}

	input := &manifest.InsertRequirementInput{
		Parent:  dest,
		Source:  opts.url,
		Version: v,
	}

	if err := opts.manifester.InsertRequirement(ctx, input); err != nil {
		return err
	}

	if err := opts.manifester.Upgrade(opts.workDir, dest); err != nil {
		return err
	}

	return nil
}
