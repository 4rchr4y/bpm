package manager

import (
	"fmt"
	"io"
	"io/fs"
	"os"

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

type (
	GetCmdHub struct {
		OsWrap      getCmdOSWrapper
		TomlEncoder getCmdTOMLEncoder
		FileLoader  getCmdLoader
	}

	GetCmdInput struct {
		Version string // specified bundle version that should be installed
		URL     string // url to the specified repository with bundle
	}

	GetCmdResult struct{}
)

type getCommand = Command[*GetCmdHub, *GetCmdInput, *GetCmdResult]

func NewGetCommand(hub *GetCmdHub) Commander {
	return &getCommand{
		Name:     GetCmdName,
		Hub:      hub,
		Run:      runGetCmd,
		Registry: NewRegistry(1),
		Requires: []string{
			InstallCmdName,
		},
	}
}

func runGetCmd(cmd *getCommand, input *GetCmdInput) (*GetCmdResult, error) {
	installCmd, err := cmd.Registry.get(InstallCmdName)
	if err != nil {
		return nil, err
	}

	inp := &InstallCmdInput{
		Version: input.Version,
		URL:     input.URL,
	}

	result, err := ExecuteInstallCmd(installCmd, inp)
	if err != nil {
		return nil, err
	}

	fmt.Println(result.Bundle.BundleFile.Package.Name)

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
