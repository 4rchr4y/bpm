package install

import (
	"github.com/4rchr4y/bpm/cli/require"
	"github.com/4rchr4y/bpm/pkg/command/factory"
	"github.com/4rchr4y/bpm/pkg/install"
	"github.com/4rchr4y/bpm/pkg/load/gitload"
	"github.com/spf13/cobra"
)

func NewCmdInstall(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install REPOSITORY",
		Aliases: []string{"i"},
		Args:    require.ExactArgs(1),
		Short:   "Install a package from the specified repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			version, err := cmd.Flags().GetString("version")
			if err != nil {
				return err
			}

			return installRun(&installOptions{
				URL:       args[0],
				Version:   version,
				GitLoader: f.GitLoader,
				Installer: f.Installer,
			})
		},
	}

	cmd.Flags().StringP("version", "v", "", "Bundle version")
	return cmd
}

type installOptions struct {
	URL       string                   // bundle repository that needs to be installed
	Version   string                   // specified bundle version
	GitLoader *gitload.GitLoader       // bundle file loader from the git repo
	Installer *install.BundleInstaller // bundle installer into the file system
}

func installRun(opts *installOptions) error {
	b, err := opts.GitLoader.DownloadBundle(opts.URL, opts.Version)
	if err != nil {
		return err
	}

	if err := opts.Installer.Install(b); err != nil {
		return err
	}

	return nil
}
