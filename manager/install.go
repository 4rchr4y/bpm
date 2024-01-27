package manager

import (
	"io/fs"
	"os"

	"github.com/4rchr4y/bpm/bundle"
	"github.com/4rchr4y/bpm/bundle/bundlefile"
)

type installCmdOSWrapper interface {
	Mkdir(name string, perm fs.FileMode) error
	Stat(name string) (fs.FileInfo, error)
	MkdirAll(path string, perm fs.FileMode) error
	Create(name string) (*os.File, error)
	UserHomeDir() (string, error)
	WriteFile(name string, data []byte, perm fs.FileMode) error
}

type installCmdLoader interface {
	DownloadBundle(url string, tag string) (*bundle.Bundle, error)
}

type installCmdBundleInstaller interface {
	Install(input *BundleInstallInput) error
}

type installCmdBfEncoder interface {
	EncodeBundleFile(bundlefile *bundlefile.File) []byte
}

type (
	InstallCmdResources struct {
		OsWrap            installCmdOSWrapper
		FileLoader        installCmdLoader
		BundleInstaller   installCmdBundleInstaller
		BundleFileEncoder initCmdBfEncoder
	}

	InstallCmdInput struct {
		Version string // specified bundle version that should be installed
		URL     string // url to the specified repository with bundle
	}

	InstallCmdResult struct {
		Bundle *bundle.Bundle
	}
)

type installCommand = Command[*InstallCmdResources, *InstallCmdInput, *InstallCmdResult]
type installCommandGuardFn = func(*installCommand, *InstallCmdInput) error

func NewInstallCommand(resources *InstallCmdResources) Commander {
	return &installCommand{
		Name:      InstallCmdName,
		Resources: resources,
		Run:       runInstallCmd,
		Guards: []installCommandGuardFn{
			validateInstallCmdInputURL,
		},
	}
}

func runInstallCmd(cmd *installCommand, input *InstallCmdInput) (*InstallCmdResult, error) {
	b, err := cmd.Resources.FileLoader.DownloadBundle(input.URL, input.Version)
	if err != nil {
		return nil, err
	}

	homeDir, err := cmd.Resources.OsWrap.UserHomeDir()
	if err != nil {
		return nil, err
	}

	installInput := &BundleInstallInput{
		HomeDir: homeDir,
		Bundle:  b,
	}

	if err := cmd.Resources.BundleInstaller.Install(installInput); err != nil {
		return nil, err
	}

	return &InstallCmdResult{
		Bundle: b,
	}, nil
}

func validateInstallCmdInputURL(cmd *installCommand, input *InstallCmdInput) error {
	return validateRepoURL(input.URL)
}
