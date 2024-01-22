package manager

import (
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/4rchr4y/bpm/bundle"
	"github.com/4rchr4y/bpm/constant"
)

type getCmdOSWrapper interface {
	Mkdir(name string, perm fs.FileMode) error
	Stat(name string) (fs.FileInfo, error)
	MkdirAll(path string, perm fs.FileMode) error
	Open(name string) (*os.File, error)
	Create(name string) (*os.File, error)
	UserHomeDir() (string, error)
	WriteFile(name string, data []byte, perm fs.FileMode) error
}

type getCmdTOMLEncoder interface {
	Encode(w io.Writer, v interface{}) error
}

type getCmdGitLoader interface {
	DownloadBundle(url string, tag string) (*bundle.Bundle, error)
}

type getCmdFsLoader interface {
	LoadBundle(dirPath string) (*bundle.Bundle, error)
}

type (
	GetCmdResources struct {
		OsWrap      getCmdOSWrapper
		TomlEncoder getCmdTOMLEncoder
		GitLoader   getCmdGitLoader
		FsLoader    getCmdFsLoader
	}

	GetCmdInput struct {
		Dir     string // bundle working directory
		Version string // specified bundle version that should be installed
		URL     string // url to the specified repository with bundle
	}

	GetCmdResult struct{}
)

type getCommand = Command[*GetCmdResources, *GetCmdInput, *GetCmdResult]
type getCommandGuardFn = func(*getCommand, *GetCmdInput) error

func NewGetCommand(resources *GetCmdResources) Commander {
	requires := []string{
		InstallCmdName,
	}

	return &getCommand{
		Name:      GetCmdName,
		Resources: resources,
		Run:       runGetCmd,
		Requires:  requires,
		Registry:  NewRegistry(len(requires)),
		Guards: []getCommandGuardFn{
			validateGetCmdBundleDir,
		},
	}
}

func runGetCmd(cmd *getCommand, input *GetCmdInput) (*GetCmdResult, error) {
	b, err := cmd.Resources.FsLoader.LoadBundle(input.Dir)
	if err != nil {
		return nil, err
	}

	installCmd, err := cmd.Registry.get(InstallCmdName)
	if err != nil {
		return nil, err
	}

	installCmdInput := &InstallCmdInput{
		Version: input.Version,
		URL:     input.URL,
	}

	result, err := ExecuteInstallCmd(installCmd, installCmdInput)
	if err != nil {
		return nil, err
	}

	b.BundleFile.Dependencies[result.Bundle.Name()] = result.Bundle.Version.Version

	bundlefilePath := fmt.Sprintf("%s/%s", input.Dir, constant.BundleFileName)
	bundlefile, err := os.OpenFile(bundlefilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}

	if err := cmd.Resources.TomlEncoder.Encode(bundlefile, b.BundleFile); err != nil {
		return nil, err
	}

	//cmd.Resources.OsWrap.WriteFile()
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

func validateGetCmdBundleDir(cmd *getCommand, input *GetCmdInput) error {
	f, err := cmd.Resources.OsWrap.Stat(fmt.Sprintf("%s/%s", input.Dir, constant.BundleFileName))
	if f == nil {
		return fmt.Errorf("cannot find '%s' providing bundle", constant.BundleFileName)
	}

	return err
}
