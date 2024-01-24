package manager

import (
	"fmt"
	"io/fs"
	"path"

	"github.com/4rchr4y/bpm/bundle"
	"github.com/4rchr4y/bpm/constant"
)

type initCmdOSWrapper interface {
	WriteFile(name string, data []byte, perm fs.FileMode) error
}

type initCmdTOMLEncoder interface {
	Encode(value interface{}) ([]byte, error)
}

type initCmdBfEncoder interface {
	EncodeBundleFile(bundlefile *bundle.BundleFile) []byte
}

type (
	InitCmdResources struct {
		OsWrap            initCmdOSWrapper
		BundleFileEncoder initCmdBfEncoder
	}

	InitCmdInput struct {
		Name   string
		Author *bundle.AuthorExpr
	}

	InitCmdResult struct{}
)

type initCommand = Command[*InitCmdResources, *InitCmdInput, *InitCmdResult]

func NewInitCommand(resources *InitCmdResources) Commander {
	return &initCommand{
		Name:      InitCmdName,
		Resources: resources,
		Run:       runInitCmd,
	}
}

func runInitCmd(cmd *initCommand, input *InitCmdInput) (*InitCmdResult, error) {
	files := map[string][]byte{
		".gitignore":            gitignoreFileContent(),
		constant.BundleFileName: bundleFileContent(cmd.Resources, input.Name, input.Author),
		constant.IgnoreFile:     bpmignoreFileContent(),
	}

	for fileName, content := range files {
		if err := cmd.Resources.OsWrap.WriteFile(fileName, []byte(content), 0644); err != nil {
			return nil, fmt.Errorf("failed to write file '%s': %v", fileName, err)
		}
	}

	return nil, nil
}

func bundleFileContent(resources *InitCmdResources, repo string, author *bundle.AuthorExpr) []byte {
	repoName := path.Base(repo)
	bundlefile := &bundle.BundleFile{
		Package: &bundle.BundleFilePackage{
			Name:        repoName,
			Author:      []string{author.String()},
			Repository:  repo,
			Description: fmt.Sprintf("Some description about '%s' bundle.", repoName),
		},
	}

	return resources.BundleFileEncoder.EncodeBundleFile(bundlefile)
}

func gitignoreFileContent() []byte {
	return []byte(`bundle.lock`)
}

func bpmignoreFileContent() []byte {
	return []byte(`.git`)
}
