package manager

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"reflect"

	"github.com/4rchr4y/bpm/bundle"
)

type getCmdOSWrapper interface {
	Mkdir(name string, perm fs.FileMode) error
	Stat(name string) (fs.FileInfo, error)
	MkdirAll(path string, perm fs.FileMode) error
	Create(name string) (*os.File, error)
	UserHomeDir() (string, error)
	WriteFile(name string, data []byte, perm fs.FileMode) error
}

type getCmdTOMLEncoder interface {
	Encode(w io.Writer, v interface{}) error
}

type getCmdLoader interface {
	DownloadBundle(url string, tag string) (*bundle.Bundle, error)
}

type getCommand struct {
	cmdName     string
	encoder     getCmdTOMLEncoder
	loader      getCmdLoader
	osWrap      getCmdOSWrapper
	subregistry *Registry
}

func (cmd *getCommand) bpmCmd()      {}
func (cmd *getCommand) Name() string { return cmd.cmdName }

func (cmd *getCommand) Requires() []string {
	return []string{
		InstallCommandName,
	}
}

func (cmd *getCommand) SetCommand(c Command) error {
	if ok := cmd.subregistry.lookup(c.Name()); ok {
		return fmt.Errorf("command '%s' in '%s' in already installed", c.Name(), cmd.cmdName)
	}

	if err := cmd.subregistry.set(c); err != nil {
		return fmt.Errorf("error occurred command '%s' in command '%s' setting: %v", c.Name(), cmd.cmdName, err)
	}

	return nil
}

type GetCmdInput struct {
	Version string // bundle package version
	URL     string // url to download bundle
}

func (cmd *getCommand) Execute(rawInput interface{}) (interface{}, error) {
	input, ok := rawInput.(*GetCmdInput)
	if !ok {
		return nil, fmt.Errorf("type '%s' is invalid input type for '%s' command", reflect.TypeOf(rawInput), cmd.cmdName)
	}

	installCmd, err := cmd.subregistry.get(InstallCommandName)
	if err != nil {
		return nil, err
	}

	result, err := installCmd.Execute(&InstallCmdInput{
		Version: input.Version,
		URL:     input.URL,
	})
	if err != nil {
		return nil, err
	}

	b, ok := result.(*bundle.Bundle)
	if !ok {
		return nil, fmt.Errorf("type '%s' is invalid type for '%s' command result", reflect.TypeOf(input).Elem().Kind().String(), InstallCommandName)
	}

	fmt.Println(b.BundleFile.Package.Name)

	// homeDir, err := cmd.osWrap.UserHomeDir()
	// if err != nil {
	// 	return nil, err
	// }

	// bundle, err := cmd.loader.LoadBundle(typedInput.BundlePath)
	// if err != nil {
	// 	return nil, err
	// }

	// bpmDirPath := fmt.Sprintf("%s/%s", homeDir, constant.BPMDir)
	// bundleName := bundle.BundleFile.Package.Name
	// bundleVersion := bundle.BundleFile.Package.Version
	// bundleVersionDir := fmt.Sprintf("%s/%s/%s", bpmDirPath, bundleName, bundleVersion)

	// if !cmd.isAlreadyInstalled(bundleVersionDir) {
	// 	fmt.Printf("Bundle '%s' with version '%s' is already installed\n", bundleName, bundleVersion)
	// 	return nil, nil
	// }

	// // creating all the directories that are necessary to save files
	// if err := cmd.osWrap.MkdirAll(bundleVersionDir, 0755); err != nil {
	// 	return nil, err
	// }

	// if err := cmd.installer.Install(&BundleInstallInput{
	// 	Dir:    bundleVersionDir,
	// 	Bundle: bundle,
	// }); err != nil {
	// 	return nil, fmt.Errorf("can't install bundle '%s': %v", bundleName, err)
	// }

	return nil, nil
}

type GetCmdConf struct {
	OsWrap      getCmdOSWrapper
	TomlEncoder getCmdTOMLEncoder
	FileLoader  getCmdLoader
}

func NewGetCommand(conf *GetCmdConf) Command {
	return &getCommand{
		cmdName:     GetCommandName,
		osWrap:      conf.OsWrap,
		encoder:     conf.TomlEncoder,
		loader:      conf.FileLoader,
		subregistry: NewRegistry(1),
	}
}
