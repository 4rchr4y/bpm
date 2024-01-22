package manager

import (
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/4rchr4y/bpm/bundle"
)

type installCmdOSWrapper interface {
	Mkdir(name string, perm fs.FileMode) error
	Stat(name string) (fs.FileInfo, error)
	MkdirAll(path string, perm fs.FileMode) error
	Create(name string) (*os.File, error)
	UserHomeDir() (string, error)
	WriteFile(name string, data []byte, perm fs.FileMode) error
}

type installCmdTOMLEncoder interface {
	Encode(w io.Writer, v interface{}) error
}

type installCmdLoader interface {
	DownloadBundle(url string, tag string) (*bundle.Bundle, error)
}

type installCmdBundleInstaller interface {
	Install(input *BundleInstallInput) error
}

type (
	InstallCmdResources struct {
		OsWrap          installCmdOSWrapper
		FileLoader      installCmdLoader
		BundleInstaller installCmdBundleInstaller
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
	repoURL := fmt.Sprintf("https://%s.git", input.URL)
	b, err := cmd.Resources.FileLoader.DownloadBundle(repoURL, input.Version)
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
