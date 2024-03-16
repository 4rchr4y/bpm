package check

import (
	"context"

	"github.com/4rchr4y/bpm/bundleutil/inspect"
	"github.com/4rchr4y/bpm/bundleutil/manifest"
	"github.com/4rchr4y/bpm/cli/cmdutil/factory"
	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/storage"
	"github.com/spf13/cobra"
)

func NewCmdTidy(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tidy [PATH]",
		Short: "Clean and inspect specified bundle",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}

			return tidyRun(cmd.Context(), &tidyOptions{
				dir:        dir,
				io:         f.IOStream,
				storage:    f.Storage,
				inspector:  f.Inspector,
				manifester: f.Manifester,
			})
		},
	}

	return cmd
}

type tidyOptions struct {
	dir        string // specified bundle folder that should be verified
	io         core.IO
	storage    *storage.Storage
	inspector  *inspect.Inspector
	manifester *manifest.Manifester
}

func tidyRun(ctx context.Context, opts *tidyOptions) error {
	b, err := opts.storage.LoadFromAbs(opts.dir, nil)
	if err != nil {
		return err
	}

	if err := opts.manifester.SyncLockfile(ctx, b); err != nil {
		return err
	}

	if err := opts.inspector.Inspect(b); err != nil {
		return err
	}

	if err := opts.manifester.Upgrade(opts.dir, b); err != nil {
		return err
	}

	opts.io.PrintfOk("bundle %s", b.BundleFile.Package.Repository)
	return nil
}
