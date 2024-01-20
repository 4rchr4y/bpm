package bpm

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"reflect"
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
	DownloadBundle(url string, tag string) (*DownloadResult, error)
}

type installCmdBundleInstaller interface {
	Install(input *BundleInstallInput) error
}

type installCommand struct {
	cmdName   string
	encoder   installCmdTOMLEncoder
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

	result, err := cmd.loader.DownloadBundle(input.URL, input.Version)
	if err != nil {
		return nil, err
	}

	homeDir, err := cmd.osWrap.UserHomeDir()
	if err != nil {
		return nil, err
	}

	installInput := &BundleInstallInput{
		HomeDir: homeDir,
		Bundle:  result.Bundle,
	}

	if err := cmd.installer.Install(installInput); err != nil {
		return nil, err
	}

	// for filePath, file := range result.Bundle.RegoFiles {
	// 	pathToSave := fmt.Sprintf(".bpm/%s/%s/%s", "test", "0.0.1", filePath)
	// 	dirToSave := filepath.Dir(pathToSave)

	// 	if _, err := os.Stat(dirToSave); os.IsNotExist(err) {
	// 		if err := os.MkdirAll(dirToSave, 0755); err != nil {
	// 			return nil, fmt.Errorf("failed to create directory '%s': %v", dirToSave, err)
	// 		}
	// 	} else if err != nil {
	// 		return nil, fmt.Errorf("error checking directory '%s': %v", dirToSave, err)
	// 	}

	// 	if err := cmd.osWrap.WriteFile(pathToSave, file.Raw, 0644); err != nil {
	// 		return nil, fmt.Errorf("failed to write file '%s': %v", pathToSave, err)
	// 	}
	// }

	return nil, nil
}

type InstallCmdConf struct {
	OsWrap          installCmdOSWrapper
	TomlEncoder     installCmdTOMLEncoder
	FileLoader      installCmdLoader
	BundleInstaller installCmdBundleInstaller
}

func NewInstallCommand(conf *InstallCmdConf) Command {
	return &installCommand{
		cmdName:   GetCommandName,
		osWrap:    conf.OsWrap,
		encoder:   conf.TomlEncoder,
		installer: conf.BundleInstaller,
		loader:    conf.FileLoader,
	}
}
