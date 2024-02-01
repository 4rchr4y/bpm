package bundleutil

import (
	"context"
	"fmt"

	"github.com/4rchr4y/bpm/pkg/bundle"
)

type downloaderGitLoader interface {
	DownloadBundle(ctx context.Context, url string, tag string) (*bundle.Bundle, error)
}

type Downloader struct {
	git downloaderGitLoader
}

func NewDownloader(git downloaderGitLoader) *Downloader {
	return &Downloader{
		git: git,
	}
}

type DownloadResult struct {
	Target    *bundle.Bundle   // target bundle that needed to be downloaded
	Rdirect   []*bundle.Bundle // directly required bundles
	Rindirect []*bundle.Bundle // indirectly required bundles
}

func (d *Downloader) Download(ctx context.Context, url string, version string) (*DownloadResult, error) {
	fmt.Println(url)
	target, err := d.git.DownloadBundle(ctx, url, version)
	if err != nil {
		return nil, err
	}

	rindirect := make([]*bundle.Bundle, 0)
	rdirect := make([]*bundle.Bundle, len(target.BundleFile.Require.List))
	for i, r := range target.BundleFile.Require.List {
		result, err := d.Download(ctx, r.Repository, r.Version)
		if err != nil {
			return nil, err
		}

		rdirect[i] = result.Target
		rindirect = append(rindirect, result.Rdirect...)
	}

	return &DownloadResult{
		Target:    target,
		Rdirect:   rdirect,
		Rindirect: rindirect,
	}, nil
}
