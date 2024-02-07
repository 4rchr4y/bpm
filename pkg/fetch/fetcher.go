package fetch

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/bpm/pkg/bundleutil"
	"github.com/4rchr4y/godevkit/syswrap/ioiface"
	"github.com/4rchr4y/godevkit/syswrap/osiface"
	"github.com/go-git/go-git/v5"
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

type Fetcher struct {
	IO     core.IO
	OSWrap osiface.OSWrapper
	IOWrap ioiface.IOWrapper

	Verifier  fetcherVerifier
	Fileifier fetcherFileifier
	GitFacade fetcherGitFacade
}

func shouldIgnore(ignoreList map[string]struct{}, path string) bool {
	if path == "" || len(ignoreList) == 0 {
		return false
	}

	dir := filepath.Dir(path)
	if dir == "." {
		return false
	}

	topLevelDir := strings.Split(dir, string(filepath.Separator))[0]
	_, found := ignoreList[topLevelDir]
	return found
}
