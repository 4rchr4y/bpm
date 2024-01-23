package manager

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

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
	Encode(value interface{}) ([]byte, error)
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

	if _, exist := b.BundleFile.Require[input.URL]; exist {
		log.Println("Already installed")
		return nil, nil
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

	if b.BundleFile.Require == nil {
		b.BundleFile.Require = make(map[string]string)
	}

	fmt.Println(b.BundleFile.Require == nil)

	b.BundleFile.Require[result.Bundle.BundleFile.Package.Repository] = getVersionStr(result.Bundle.Version)

	bundlefilePath := filepath.Join(input.Dir, constant.BundleFileName)
	// bundlefile, err := os.OpenFile(bundlefilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	// if err != nil {
	// 	return nil, err
	// }

	bytes, err := cmd.Resources.TomlEncoder.Encode(b.BundleFile)
	if err != nil {
		return nil, err
	}

	if err := cmd.Resources.OsWrap.WriteFile(bundlefilePath, bytes, 0644); err != nil {
		return nil, err
	}

	return nil, nil
}

func validateGetCmdBundleDir(cmd *getCommand, input *GetCmdInput) error {
	f, err := cmd.Resources.OsWrap.Stat(filepath.Join(input.Dir, constant.BundleFileName))
	if f == nil {
		return fmt.Errorf("cannot find '%s' providing bundle", constant.BundleFileName)
	}

	return err
}

func getVersionStr(version *bundle.VersionExpr) string {
	if version.Version != constant.BundlePseudoVersion {
		return version.Version
	}

	return version.String()
}
