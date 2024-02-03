package bundleutil

import (
	"context"
	"fmt"
	"sync"

	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/bpm/pkg/bundle/bundlefile"
)

type downloaderGitLoader interface {
	DownloadBundle(ctx context.Context, url string, tag string) (*bundle.Bundle, error)
}

type downloaderVerifier interface {
	Verify(b *bundle.Bundle) error
}

type Downloader struct {
	verifier downloaderVerifier
	git      downloaderGitLoader
}

func NewDownloader(git downloaderGitLoader, verifier downloaderVerifier) *Downloader {
	return &Downloader{
		git:      git,
		verifier: verifier,
	}
}

type DownloadResult struct {
	Target    *bundle.Bundle   // target bundle that needed to be downloaded
	Rdirect   []*bundle.Bundle // directly required bundles
	Rindirect []*bundle.Bundle // indirectly required bundles
}

func (dr *DownloadResult) Merge() []*bundle.Bundle {
	totalLen := len(dr.Rdirect) + len(dr.Rindirect)
	if dr.Target != nil {
		totalLen++
	}

	result := make([]*bundle.Bundle, totalLen)

	index := 0
	if dr.Target != nil {
		result[index] = dr.Target
		index++
	}

	index += copy(result[index:], dr.Rdirect)
	copy(result[index:], dr.Rindirect)

	return result
}

func (d *Downloader) Download(ctx context.Context, url string, version string) (*DownloadResult, error) {
	target, err := d.git.DownloadBundle(ctx, url, version)
	if err != nil {
		return nil, fmt.Errorf("failed to download '%s' bundle: %w", url, err)
	}

	// if err := d.verifier.Verify(target); err != nil {
	// 	return nil, err
	// }

	if target.BundleFile.Require == nil {
		return &DownloadResult{Target: target}, nil
	}

	result, err := d.downloadRequires(ctx, target.BundleFile.Require.List)
	if err != nil {
		return nil, err
	}

	return &DownloadResult{
		Target:    target,
		Rdirect:   result.direct,
		Rindirect: result.indirect,
	}, nil
}

type downloadRequiresResults struct {
	direct   []*bundle.Bundle
	indirect []*bundle.Bundle
}

func (d *Downloader) downloadRequires(ctx context.Context, requires []*bundlefile.RequirementDecl) (*downloadRequiresResults, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	resultsChan := make(chan *DownloadResult)
	errChan := make(chan error)

	for _, req := range requires {
		wg.Add(1)
		go func(req *bundlefile.RequirementDecl) {
			defer wg.Done()
			result, err := d.Download(ctx, req.Repository, req.Version)
			if err != nil {
				select {
				case errChan <- fmt.Errorf("failed to download required '%s': %w", req.Repository, err):
				case <-ctx.Done():
				}
				return
			}
			select {
			case resultsChan <- result:
			case <-ctx.Done():
			}
		}(req)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
		close(errChan)
	}()

	result := new(downloadRequiresResults)
	for i := 0; i < len(requires); i++ {
		select {
		case err := <-errChan:
			return nil, err
		case r := <-resultsChan:
			result.direct = append(result.direct, r.Target)
			result.indirect = append(result.indirect, r.Rdirect...)
		case <-ctx.Done():
			return result, ctx.Err()
		}
	}

	return result, nil
}
