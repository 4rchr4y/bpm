package manager

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"reflect"

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

type installCommand struct {
	cmdName   string
	installer installCmdBundleInstaller
	loader    installCmdLoader
	osWrap    installCmdOSWrapper
}

func (cmd *installCommand) bpmCmd()                  {}
func (cmd *installCommand) Name() string             { return cmd.cmdName }
func (cmd *installCommand) Requires() []string       { return nil }
func (cmd *installCommand) SetCommand(Command) error { return nil }

type InstallCmdInput struct {
	Version string // bundle package version
	URL     string // url to download bundle
}

func (cmd *installCommand) Execute(rawInput interface{}) (interface{}, error) {
	input, ok := rawInput.(*InstallCmdInput)
	if !ok {
		return nil, fmt.Errorf("type '%s' is invalid input type for '%s' command", reflect.TypeOf(rawInput), cmd.cmdName)
	}

	b, err := cmd.loader.DownloadBundle(input.URL, input.Version)
	if err != nil {
		return nil, err
	}

	homeDir, err := cmd.osWrap.UserHomeDir()
	if err != nil {
		return nil, err
	}

	installInput := &BundleInstallInput{
		HomeDir: homeDir,
		Bundle:  b,
	}

	if err := cmd.installer.Install(installInput); err != nil {
		return nil, err
	}

	return b, nil
}

type InstallCmdConf struct {
	OsWrap          installCmdOSWrapper
	FileLoader      installCmdLoader
	BundleInstaller installCmdBundleInstaller
}

func NewInstallCommand(conf *InstallCmdConf) Command {
	return &installCommand{
		cmdName:   InstallCommandName,
		osWrap:    conf.OsWrap,
		installer: conf.BundleInstaller,
		loader:    conf.FileLoader,
	}
}
