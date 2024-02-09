package verify

import (
	"context"

	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/pkg/bundleutil"
	"github.com/4rchr4y/bpm/pkg/bundleutil/inspect"
	"github.com/4rchr4y/bpm/pkg/cmdutil/factory"
	"github.com/4rchr4y/bpm/pkg/cmdutil/require"
	"github.com/4rchr4y/bpm/pkg/fetch"
	"github.com/spf13/cobra"
)

func NewCmdVerify(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify PATH",
		Short: "Verify specified bundle",
		Args:  require.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return verifyRun(cmd.Context(), &verifyOptions{
				dir:        args[0],
				io:         f.IOStream,
				Fetcher:    f.Fetcher,
				inspector:  f.Inspector,
				manifester: f.Manifester,
			})
		},
	}

	return cmd
}

type verifyOptions struct {
	dir        string // specified bundle folder that should be verified
	io         core.IO
	Fetcher    *fetch.Fetcher
	inspector  *inspect.Inspector
	manifester *bundleutil.Manifester
}

func verifyRun(ctx context.Context, opts *verifyOptions) error {
	b, err := opts.Fetcher.FetchLocal(opts.dir)
	if err != nil {
		return err
	}

	if err := opts.manifester.Update(&bundleutil.UpdateInput{Target: b}); err != nil {
		return err
	}

	if err := opts.manifester.Upgrade(opts.dir, b); err != nil {
		return err
	}

	if err := opts.inspector.Verify(b); err != nil {
		return err
	}

	opts.io.PrintfOk("bundle '%s' is verified", b.Repository())
	return nil
}
