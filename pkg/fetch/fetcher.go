package fetch

import (
	"context"

	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/godevkit/v3/syswrap/ioiface"
	"github.com/4rchr4y/godevkit/v3/syswrap/osiface"
)

type fetcherStorage interface {
	Lookup(repo string, version string) bool
	Load(repo string, version *bundle.VersionExpr) (*bundle.Bundle, error)
	MakeBundleSourcePath(repo string, version string) string
}

type fetcherDownloader interface {
	PlainDownloadWithContext(ctx context.Context, url string, tag *bundle.VersionExpr) (*bundle.Bundle, error)
}

type Fetcher struct {
	IO     core.IO
	OSWrap osiface.OSWrapper
	IOWrap ioiface.IOWrapper

	Storage    fetcherStorage
	Downloader fetcherDownloader
}

type SourceType int

const (
	Remote SourceType = iota
	Local
)

var sourceTypeStr = [...]string{
	Remote: "remote",
	Local:  "local",
}

func (st SourceType) String() string { return sourceTypeStr[st] }

type FetchOutput struct {
	Source SourceType
	Bundle *bundle.Bundle
}

func (f *Fetcher) Fetch(ctx context.Context, repo string, version *bundle.VersionExpr) (*bundle.Bundle, error) {
	if b, _ := f.FetchLocal(ctx, repo, version); b != nil {
		return b, nil
	}

	b, err := f.FetchRemote(ctx, repo, version)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (f *Fetcher) FetchLocal(ctx context.Context, repo string, version *bundle.VersionExpr) (*bundle.Bundle, error) {
	if ok := f.Storage.Lookup(repo, version.String()); !ok {
		return nil, nil
	}

	b, err := f.Storage.Load(repo, version)
	if err != nil {
		f.IO.PrintfErr("failed to load bundle %s from local storage: %v", f.Storage.MakeBundleSourcePath(repo, version.String()), err)
		return nil, err
	}

	return b, nil

}

func (f *Fetcher) FetchRemote(ctx context.Context, repo string, version *bundle.VersionExpr) (*bundle.Bundle, error) {
	return f.Downloader.PlainDownloadWithContext(ctx, repo, version)
}
