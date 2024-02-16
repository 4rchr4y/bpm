package check

import (
	"context"

	"github.com/4rchr4y/bpm/bundleutil/inspect"
	"github.com/4rchr4y/bpm/bundleutil/manifest"
	"github.com/4rchr4y/bpm/cmdutil/factory"
	"github.com/4rchr4y/bpm/cmdutil/require"
	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/storage"
	"github.com/spf13/cobra"
)

func NewCmdCheck(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check PATH",
		Short: "Check specified bundle",
		Args:  require.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return checkRun(cmd.Context(), &checkOptions{
				dir:        args[0],
				io:         f.IOStream,
				Storage:    f.Storage,
				inspector:  f.Inspector,
				manifester: f.Manifester,
			})
		},
	}

	return cmd
}

type checkOptions struct {
	dir        string // specified bundle folder that should be verified
	io         core.IO
	Storage    *storage.Storage
	inspector  *inspect.Inspector
	manifester *manifest.Manifester
}

func checkRun(ctx context.Context, opts *checkOptions) error {
	b, err := opts.Storage.LoadFromAbs(opts.dir, nil)
	if err != nil {
		return err
	}

	if err := opts.manifester.CreateRequirement(&manifest.CreateRequirementInput{Parent: b}); err != nil {
		return err
	}

	if err := opts.manifester.Upgrade(opts.dir, b); err != nil {
		return err
	}

	if err := opts.inspector.Inspect(b); err != nil {
		return err
	}

	opts.io.PrintfOk("bundle '%s' is checked", b.Repository())
	return nil
}
