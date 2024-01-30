package get

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/4rchr4y/bpm/cli/require"
	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/bpm/pkg/bundle/lockfile"
	"github.com/4rchr4y/bpm/pkg/command/factory"
	"github.com/4rchr4y/bpm/pkg/encode"
	"github.com/4rchr4y/bpm/pkg/install"
	"github.com/4rchr4y/bpm/pkg/load/gitload"
	"github.com/4rchr4y/bpm/pkg/load/osload"
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

			return getRun(&getOptions{
				WorkDir:   wd,
				URL:       args[0],
				Version:   version,
				GitLoader: f.GitLoader,
				OsLoader:  f.OsLoader,
				Installer: f.Installer,
				Encoder:   f.Encoder,
				WriteFile: f.OS.WriteFile,
			})
		},
	}

	cmd.Flags().StringP("version", "v", "", "Bundle version")
	return cmd
}

type getOptions struct {
	WorkDir   string                                                 // bundle working directory
	URL       string                                                 // bundle repository that needs to be installed
	Version   string                                                 // specified bundle version
	GitLoader *gitload.GitLoader                                     // bundle file loader from the git repo
	OsLoader  *osload.OsLoader                                       // bundle file loader from file system
	Installer *install.BundleInstaller                               // bundle installer into the file system
	Encoder   *encode.BundleEncoder                                  // decoder of bundle component files
	WriteFile func(name string, data []byte, perm fs.FileMode) error // func of saving a file to disk
}

func getRun(opts *getOptions) error {
	workingBundle, err := opts.OsLoader.LoadBundle(opts.WorkDir)
	if err != nil {
		return err
	}

	b, err := opts.GitLoader.DownloadBundle(opts.URL, opts.Version)
	if err != nil {
		return err
	}

	// TODO: need to check the bundle here

	require := make([]*bundle.Bundle, len(b.BundleFile.Require.List))
	for i, r := range b.BundleFile.Require.List {
		requireBundle, err := opts.GitLoader.DownloadBundle(r.Repository, r.Version)
		if err != nil {
			return err
		}

		require[i] = requireBundle
	}

	if err := opts.Installer.Install(b); err != nil {
		return err
	}

	if err := workingBundle.SetRequire(b, lockfile.Direct); err != nil {
		return err
	}

	if err := updateBundleFile(opts, workingBundle); err != nil {
		return err
	}

	if err := updateLockFile(opts, workingBundle); err != nil {
		return err
	}

	return nil
}

func updateBundleFile(opts *getOptions, workingBundle *bundle.Bundle) error {
	bundlefilePath := filepath.Join(opts.WorkDir, constant.BundleFileName)
	bytes := opts.Encoder.EncodeBundleFile(workingBundle.BundleFile)

	if err := opts.WriteFile(bundlefilePath, bytes, 0644); err != nil {
		return fmt.Errorf("error occurred while '%s' file updating: %v", constant.BundleFileName, err)
	}

	return nil
}

func updateLockFile(opts *getOptions, workingBundle *bundle.Bundle) error {
	bundlefilePath := filepath.Join(opts.WorkDir, constant.LockFileName)
	bytes := opts.Encoder.EncodeLockFile(workingBundle.LockFile)

	if err := opts.WriteFile(bundlefilePath, bytes, 0644); err != nil {
		return fmt.Errorf("error occurred while '%s' file updating: %v", constant.LockFileName, err)
	}

	return nil
}
