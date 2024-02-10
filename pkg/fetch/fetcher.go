package fetch

import (
	"context"

	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/godevkit/v3/syswrap/ioiface"
	"github.com/4rchr4y/godevkit/v3/syswrap/osiface"
)

type fetcherStorage interface {
	Load(repo string, version *bundle.VersionExpr) (*bundle.Bundle, error)
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

func (f *Fetcher) FetchLocal(ctx context.Context, repo string, version *bundle.VersionExpr) (*bundle.Bundle, error) {
	return f.Storage.Load(repo, version)
}

func (f *Fetcher) FetchRemote(ctx context.Context, repo string, version *bundle.VersionExpr) (*bundle.Bundle, error) {
	return f.Downloader.PlainDownloadWithContext(ctx, repo, version)
}
