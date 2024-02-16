package fetch

import (
	"context"
	"fmt"

	"github.com/4rchr4y/bpm/bundle"
	"github.com/4rchr4y/bpm/core"
)

type fetcherInspector interface {
	Inspect(b *bundle.Bundle) error
}

type fetcherStorage interface {
	Some(source string, version string) bool
	Load(source string, version *bundle.VersionExpr) (*bundle.Bundle, error)
	MakeBundleSourcePath(source string, version string) string
}

type fetcherGitHub interface {
	Download(ctx context.Context, source string, tag *bundle.VersionExpr) (*bundle.Bundle, error)
}

type Fetcher struct {
	IO        core.IO
	Storage   fetcherStorage
	Inspector fetcherInspector
	GitHub    fetcherGitHub
}

type FetchResult struct {
	Target    *bundle.Bundle   // target bundle that needed to be downloaded
	Rdirect   []*bundle.Bundle // directly required bundles
	Rindirect []*bundle.Bundle // indirectly required bundles
}

func (fres *FetchResult) Merge() []*bundle.Bundle {
	totalLen := len(fres.Rdirect) + len(fres.Rindirect)
	if fres.Target != nil {
		totalLen++
	}

	result := make([]*bundle.Bundle, totalLen)

	index := 0
	if fres.Target != nil {
		result[index] = fres.Target
		index++
	}

	index += copy(result[index:], fres.Rdirect)
	copy(result[index:], fres.Rindirect)

	return result
}

func (d *Fetcher) Fetch(ctx context.Context, source string, version *bundle.VersionExpr) (*FetchResult, error) {
	target, err := d.PlainFetch(ctx, source, version)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %v", source, err)
	}

	if target.BundleFile.Require == nil {
		return &FetchResult{Target: target}, nil
	}

	rindirect := make([]*bundle.Bundle, 0)
	rdirect := make([]*bundle.Bundle, len(target.BundleFile.Require.List))
	for i, r := range target.BundleFile.Require.List {
		v, err := bundle.ParseVersionExpr(r.Version)
		if err != nil {
			return nil, err
		}

		result, err := d.Fetch(ctx, r.Repository, v)
		if err != nil {
			return nil, err
		}

		rdirect[i] = result.Target
		rindirect = append(rindirect, result.Rdirect...)
	}

	return &FetchResult{
		Target:    target,
		Rdirect:   rdirect,
		Rindirect: rindirect,
	}, nil
}

func (f *Fetcher) PlainFetch(ctx context.Context, source string, version *bundle.VersionExpr) (*bundle.Bundle, error) {
	if b, _ := f.FetchLocal(ctx, source, version); b != nil {
		return b, nil
	}

	b, err := f.FetchRemote(ctx, source, version)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (f *Fetcher) FetchLocal(ctx context.Context, source string, version *bundle.VersionExpr) (*bundle.Bundle, error) {
	ok := f.Storage.Some(source, version.String())
	if !ok {
		return nil, nil
	}

	b, err := f.Storage.Load(source, version)
	if err != nil {
		f.IO.PrintfErr("failed to load bundle %s from local storage: %v", f.Storage.MakeBundleSourcePath(source, version.String()), err)
		return nil, err
	}

	if err := f.Inspector.Inspect(b); err != nil {
		return nil, err
	}

	return b, nil
}

func (f *Fetcher) FetchRemote(ctx context.Context, source string, version *bundle.VersionExpr) (*bundle.Bundle, error) {
	b, err := f.GitHub.Download(ctx, source, version)
	if err != nil {
		return nil, err
	}

	if err := f.Inspector.Inspect(b); err != nil {
		return nil, err
	}

	return b, nil
}
