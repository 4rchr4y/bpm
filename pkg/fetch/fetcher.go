package fetch

import (
	"context"

	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/bpm/pkg/bundle/bundlefile"
	"github.com/4rchr4y/bpm/pkg/bundle/lockfile"
	"github.com/4rchr4y/bpm/pkg/bundleutil/bundlebuild"
	"github.com/4rchr4y/godevkit/syswrap/ioiface"
	"github.com/4rchr4y/godevkit/syswrap/osiface"
	"github.com/go-git/go-git/v5"
)

type fetcherGitFacade interface {
	CloneWithContext(ctx context.Context, opts *git.CloneOptions) (*git.Repository, error)
}

type fetcherInspector interface {
	Inspect(b *bundle.Bundle) error
}

type fetcherEncoder interface {
	DecodeBundleFile(content []byte) (*bundlefile.File, error)
	DecodeLockFile(content []byte) (*lockfile.File, error)
	DecodeIgnoreFile(content []byte) (*bundle.IgnoreFile, error)
	Fileify(files map[string][]byte, options ...bundlebuild.BundleOptFn) (*bundle.Bundle, error)
}

type fetcherStorage interface {
	Load(dirPath string) (*bundle.Bundle, error)
}

type Fetcher struct {
	IO     core.IO
	OSWrap osiface.OSWrapper
	IOWrap ioiface.IOWrapper

	Storage   fetcherStorage
	Inspector fetcherInspector
	GitFacade fetcherGitFacade
	Encoder   fetcherEncoder
}
