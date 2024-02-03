package verify

import (
	"context"
	"fmt"

	"github.com/4rchr4y/bpm/cli/require"
	"github.com/4rchr4y/bpm/pkg/bundleutil"
	"github.com/4rchr4y/bpm/pkg/cmdutil/factory"
	"github.com/4rchr4y/bpm/pkg/load/osload"
	"github.com/spf13/cobra"
)

func NewCmdVerify(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify PATH",
		Short: "Verify specified bundle",
		Args:  require.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return verifyRun(cmd.Context(), &verifyOptions{
				Dir:      args[0],
				OsLoader: f.OsLoader,
				Verifier: f.Verifier,
			})
		},
	}

	return cmd
}

type verifyOptions struct {
	Dir      string           // specified bundle folder that should be verified
	OsLoader *osload.OsLoader // bundle file loader from file system
	Verifier *bundleutil.Verifier
}

func verifyRun(ctx context.Context, opts *verifyOptions) error {
	b, err := opts.OsLoader.LoadBundle(opts.Dir)
	if err != nil {
		return err
	}

	result, err := opts.Verifier.Verify(b)
	fmt.Println(result)

	return err
}
