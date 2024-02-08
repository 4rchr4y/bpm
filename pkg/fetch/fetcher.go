package fetch

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/bpm/pkg/bundle/bundlefile"
	"github.com/4rchr4y/bpm/pkg/bundle/lockfile"
	"github.com/4rchr4y/bpm/pkg/bundle/regofile"
	"github.com/4rchr4y/bpm/pkg/bundleutil"
	"github.com/4rchr4y/godevkit/syswrap/ioiface"
	"github.com/4rchr4y/godevkit/syswrap/osiface"
	"github.com/go-git/go-git/v5"
	"github.com/open-policy-agent/opa/ast"
)

type fetcherGitFacade interface {
	CloneWithContext(ctx context.Context, opts *git.CloneOptions) (*git.Repository, error)
}

type fetcherFileifier interface {
	Fileify(files map[string][]byte, options ...bundleutil.BundleOptFn) (*bundle.Bundle, error)
}

type fetcherVerifier interface {
	Verify(b *bundle.Bundle) error
}

type fetcherEncoder interface {
	DecodeBundleFile(content []byte) (*bundlefile.File, error)
	DecodeLockFile(content []byte) (*lockfile.File, error)
	DecodeIgnoreFile(content []byte) (*bundle.IgnoreFile, error)
}

type fetcherStorage interface {
	Load(dirPath string) (*bundle.Bundle, error)
}

type Fetcher struct {
	IO     core.IO
	OSWrap osiface.OSWrapper
	IOWrap ioiface.IOWrapper

	Storage fetcherStorage

	Verifier  fetcherVerifier
	GitFacade fetcherGitFacade
	Encoder   fetcherEncoder
}

// transformer of file contents into structures
func (fetcher *Fetcher) Fileify(files map[string][]byte, options ...bundleutil.BundleOptFn) (*bundle.Bundle, error) {
	b := &bundle.Bundle{
		RegoFiles:  make(map[string]*regofile.File),
		OtherFiles: make(map[string][]byte),
	}

	for i := range options {
		options[i](b)
	}

	modules := &lockfile.ModulesDecl{
		List: make([]*lockfile.ModDecl, 0),
	}

	for filePath, content := range files {
		switch {
		case isRegoFile(filePath):
			parsed, err := ast.ParseModule(filePath, string(content))
			if err != nil {
				return nil, fmt.Errorf("error parsing file contents: %v", err)
			}

			file := &regofile.File{Path: filePath, Parsed: parsed}
			b.RegoFiles[filePath] = file

			modules.List = append(modules.List, &lockfile.ModDecl{
				Package: file.Package(),
				Source:  filePath,
				Sum:     file.Sum(),
			})

		case isBPMFile(filePath):
			bundlefile, err := fetcher.Encoder.DecodeBundleFile(content)
			if err != nil {
				return nil, fmt.Errorf("error occurred while decoding %s content: %v", constant.BundleFileName, err)
			}

			b.BundleFile = bundlefile

		case isBPMLockFile(filePath):
			lockfile, err := fetcher.Encoder.DecodeLockFile(content)
			if err != nil {
				return nil, fmt.Errorf("error occurred while decoding %s content: %v", constant.BundleFileName, err)
			}

			b.LockFile = lockfile

		default:
			b.OtherFiles[filePath] = content
		}
	}

	if b.LockFile == nil {
		fetcher.IO.PrintfWarn("file '%s' in bundle '%s' was not found", constant.LockFileName, b.Name())

		if err := initLockFile(b); err != nil {
			return nil, err
		}
	}

	if len(modules.List) > 0 {
		b.LockFile.Modules = modules.Sort()
	}

	return b, nil
}

func isRegoFile(filePath string) bool    { return filepath.Ext(filePath) == constant.RegoFileExt }
func isBPMFile(filePath string) bool     { return filePath == constant.BundleFileName }
func isBPMLockFile(filePath string) bool { return filePath == constant.LockFileName }
func isEmpty(content []byte) bool        { return len(content) < 1 }

func initLockFile(b *bundle.Bundle) error {
	if b.BundleFile == nil {
		return fmt.Errorf("can't find '%s' file", constant.BundleFileName)
	}

	b.LockFile = &lockfile.File{
		// TODO: set 'edition' from global app context
		Edition: "2024",
		Sum:     b.Sum(),
	}

	return nil
}
