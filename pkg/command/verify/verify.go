package verify

import (
	"context"

	"github.com/4rchr4y/bpm/pkg/bundleutil"
	"github.com/4rchr4y/bpm/pkg/cmdutil/factory"
	"github.com/4rchr4y/bpm/pkg/cmdutil/require"
	"github.com/spf13/cobra"
)

func NewCmdVerify(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify PATH",
		Short: "Verify specified bundle",
		Args:  require.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return verifyRun(cmd.Context(), &verifyOptions{
				Dir:        args[0],
				Loader:     f.Loader,
				Verifier:   f.Verifier,
				Manifester: f.Manifester,
			})
		},
	}

	return cmd
}

type verifyOptions struct {
	Dir        string // specified bundle folder that should be verified
	Loader     *bundleutil.Loader
	Verifier   *bundleutil.Verifier
	Manifester *bundleutil.Manifester
}

func verifyRun(ctx context.Context, opts *verifyOptions) error {
	b, err := opts.Loader.LoadBundle(opts.Dir)
	if err != nil {
		return err
	}

	currentChecksum := b.Sum()
	if currentChecksum != b.LockFile.Sum {
		b.LockFile.Sum = currentChecksum
	}

	if err := opts.Manifester.Upgrade(opts.Dir, b); err != nil {
		return err
	}

	return opts.Verifier.Verify(b)
}
