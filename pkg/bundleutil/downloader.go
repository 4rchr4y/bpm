package bundleutil

import (
	"context"
	"fmt"

	"github.com/4rchr4y/bpm/pkg/bundle"
)

type downloaderGitLoader interface {
	FetchRemote(ctx context.Context, url string, tag *bundle.VersionExpr) (*bundle.Bundle, error)
}

type downloaderInspector interface {
	Inspect(b *bundle.Bundle) error
}

type Downloader struct {
	inspector downloaderInspector
	git       downloaderGitLoader
}

func NewDownloader(git downloaderGitLoader, verifier downloaderInspector) *Downloader {
	return &Downloader{
		git:       git,
		inspector: verifier,
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

func (d *Downloader) Download(ctx context.Context, url string, version *bundle.VersionExpr) (*DownloadResult, error) {
	target, err := d.git.FetchRemote(ctx, url, version)
	if err != nil {
		return nil, fmt.Errorf("failed to download '%s' bundle: %w", url, err)
	}

	// TODO: probably should be moved to Loader
	if err := d.inspector.Inspect(target); err != nil {
		return nil, err
	}

	if target.BundleFile.Require == nil {
		return &DownloadResult{Target: target}, nil
	}

	rindirect := make([]*bundle.Bundle, 0)
	rdirect := make([]*bundle.Bundle, len(target.BundleFile.Require.List))
	for i, r := range target.BundleFile.Require.List {
		v, err := bundle.ParseVersionExpr(r.Version)
		if err != nil {
			return nil, err
		}

		result, err := d.Download(ctx, r.Repository, v)
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

// func (d *Downloader) Download(ctx context.Context, url string, version *bundle.VersionExpr) (*DownloadResult, error) {
// 	// 1. Fetch local bundle storage

// 	return nil, nil
// }

func (d *Downloader) download(ctx context.Context, url string, v *bundle.VersionExpr) (*DownloadResult, error) {
	return nil, nil
}

// func (d *Downloader) Download(ctx context.Context, url string, version string) (*DownloadResult, error) {
// 	target, err := d.git.DownloadBundle(ctx, url, version)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to download '%s' bundle: %w", url, err)
// 	}

// 	if err := d.verifier.Verify(target); err != nil {
// 		return nil, err
// 	}

// 	if target.BundleFile.Require == nil {
// 		return &DownloadResult{Target: target}, nil
// 	}

// 	result, err := d.downloadRequires(ctx, target.BundleFile.Require.List)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &DownloadResult{
// 		Target:    target,
// 		Rdirect:   result.direct,
// 		Rindirect: result.indirect,
// 	}, nil
// }

// type downloadRequiresResults struct {
// 	direct   []*bundle.Bundle
// 	indirect []*bundle.Bundle
// }

// func (d *Downloader) downloadRequires(ctx context.Context, requires []*bundlefile.RequirementDecl) (*downloadRequiresResults, error) {
// 	ctx, cancel := context.WithCancel(ctx)
// 	defer cancel()

// 	var wg sync.WaitGroup
// 	resultsChan := make(chan *DownloadResult)
// 	errChan := make(chan error)

// 	for _, req := range requires {
// 		wg.Add(1)
// 		go func(req *bundlefile.RequirementDecl) {
// 			defer wg.Done()
// 			result, err := d.Download(ctx, req.Repository, req.Version)
// 			if err != nil {
// 				select {
// 				case errChan <- fmt.Errorf("failed to download required '%s': %w", req.Repository, err):
// 				case <-ctx.Done():
// 				}
// 				return
// 			}
// 			select {
// 			case resultsChan <- result:
// 			case <-ctx.Done():
// 			}
// 		}(req)
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(resultsChan)
// 		close(errChan)
// 	}()

// 	result := new(downloadRequiresResults)
// 	for i := 0; i < len(requires); i++ {
// 		select {
// 		case err := <-errChan:
// 			return nil, err
// 		case r := <-resultsChan:
// 			result.direct = append(result.direct, r.Target)
// 			result.indirect = append(result.indirect, r.Rdirect...)
// 		case <-ctx.Done():
// 			return result, ctx.Err()
// 		}
// 	}

// 	return result, nil
// }
