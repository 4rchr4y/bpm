package fetch

import (
	"context"
	"fmt"
	"strings"

	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/godevkit/v3/must"
	"github.com/4rchr4y/godevkit/v3/syswrap/ioiface"
	"github.com/4rchr4y/godevkit/v3/syswrap/osiface"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hashicorp/go-version"
)

type fetcherGitFacade interface {
	CloneWithContext(ctx context.Context, opts *git.CloneOptions) (*git.Repository, error)
}

type fetcherInspector interface {
	Inspect(b *bundle.Bundle) error
}

type fetcherEncoder interface {
	DecodeIgnoreFile(content []byte) (*bundle.IgnoreFile, error)
	Fileify(files map[string][]byte) (*bundle.BundleRaw, error)
}

type fetcherStorage interface {
	Lookup(source string, version string) bool
	Load(source string, version *bundle.VersionExpr) (*bundle.Bundle, error)
	MakeBundleSourcePath(source string, version string) string
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

	// if err := d.Inspector.Inspect(target); err != nil {
	// 	return nil, err
	// }

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
	ok := f.Storage.Lookup(source, version.String())
	if !ok {
		return nil, nil
	}
	fmt.Println(ok)
	b, err := f.Storage.Load(source, version)
	if err != nil {
		f.IO.PrintfErr("failed to load bundle %s from local storage: %v", f.Storage.MakeBundleSourcePath(source, version.String()), err)
		return nil, err
	}

	return b, nil

}

func (f *Fetcher) FetchRemote(ctx context.Context, source string, version *bundle.VersionExpr) (*bundle.Bundle, error) {
	return f.Download(ctx, source, version)
}

func (d *Fetcher) Download(ctx context.Context, source string, tag *bundle.VersionExpr) (*bundle.Bundle, error) {
	d.IO.PrintfInfo("downloading %s@%s", source, tag.String())

	options := &git.CloneOptions{
		URL: fmt.Sprintf("https://%s.git", source),
	}

	repo, err := d.GitFacade.CloneWithContext(ctx, options)
	if err != nil {
		return nil, err
	}

	commit, v, err := fetchCommitByTag(repo, tag)
	if err != nil {
		return nil, err
	}

	files, ignoreFile, err := d.getFilesFromCommit(commit)
	if err != nil {
		return nil, err
	}

	bundleRaw, err := d.Encoder.Fileify(files)
	if err != nil {
		return nil, err
	}

	return bundleRaw.ToBundle(v, ignoreFile)
}

func fetchCommitByTag(repo *git.Repository, v *bundle.VersionExpr) (*object.Commit, *bundle.VersionExpr, error) {
	if v == nil {
		return getLatestVersionCommit(repo)
	}

	if v.IsPseudo() {
		commit, err := getPseudoVersionCommit(repo, v)
		if err != nil {
			return nil, nil, err
		}

		return commit, v, nil
	}

	commit, err := getCurrentVersionCommit(repo, v.Tag.Original())
	if err != nil {
		return nil, nil, err
	}

	return commit, v, nil

}

func getPseudoVersionCommit(repo *git.Repository, v *bundle.VersionExpr) (*object.Commit, error) {
	if len(v.Hash) != bundle.VersionShortHashLen {
		return nil, fmt.Errorf("short hash must be %d characters long", bundle.VersionShortHashLen)
	}

	iter, err := repo.CommitObjects()
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var commit *object.Commit
	err = iter.ForEach(func(c *object.Commit) error {
		if strings.HasPrefix(c.Hash.String(), v.Hash) && v.Timestamp == c.Committer.When.UTC() {
			commit = c
			return nil
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if commit == nil {
		return nil, fmt.Errorf("commit not found with hash '%s'", v.Hash)
	}

	return commit, nil
}

func getLatestVersionCommit(repo *git.Repository) (*object.Commit, *bundle.VersionExpr, error) {
	tags, err := collectTagList(repo)
	if err != nil {
		return nil, nil, err
	}

	v, ref := mustFindLatestVersion(tags)
	if v == nil || ref == nil {
		commit, err := getLatestCommit(repo)
		if err != nil {
			return nil, nil, err
		}

		return commit, bundle.NewVersionExprFromCommit(commit, v), nil
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, nil, err
	}

	return commit, bundle.NewVersionExprFromCommit(commit, v), nil
}

func collectTagList(repo *git.Repository) (map[*version.Version]*plumbing.Reference, error) {
	iter, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	tags := make(map[*version.Version]*plumbing.Reference)
	err = iter.ForEach(func(ref *plumbing.Reference) error {
		if !ref.Name().IsTag() {
			return nil
		}

		v, err := version.NewVersion(ref.Name().Short())
		if err != nil {
			return nil
		}

		tags[v] = ref

		return nil
	})
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func mustFindLatestVersion(tags map[*version.Version]*plumbing.Reference) (v *version.Version, ref *plumbing.Reference) {
	// use MUST here because the constant version should always be correct
	v = must.Must(version.NewSemver(constant.BundlePseudoVersion))

	for version, reference := range tags {
		if version.GreaterThan(v) {
			v = version
			ref = reference
		}
	}

	return v, ref
}

func getLatestCommit(repo *git.Repository) (*object.Commit, error) {
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}
	return repo.CommitObject(ref.Hash())
}

func getCurrentVersionCommit(repo *git.Repository, tag string) (*object.Commit, error) {
	tags, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	var ref *plumbing.Reference
	err = tags.ForEach(func(r *plumbing.Reference) error {
		if !r.Name().IsTag() {
			return nil
		}

		if r.Name().Short() == tag {
			ref = r
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if ref == nil {
		return nil, fmt.Errorf("version '%s' is not found", tag)
	}

	return repo.CommitObject(ref.Hash())
}

func (d *Fetcher) getFilesFromCommit(commit *object.Commit) (map[string][]byte, *bundle.IgnoreFile, error) {
	filesIter, err := commit.Files()
	if err != nil {
		return nil, nil, err
	}

	ignoreFile, err := d.fetchIgnoreListFromGitCommit(commit)
	if err != nil {
		return nil, nil, err
	}

	files := make(map[string][]byte)
	err = filesIter.ForEach(func(f *object.File) error {
		if ignoreFile.Some(f.Name) {
			return nil
		}

		content, err := f.Contents()
		if err != nil {
			return err
		}

		files[f.Name] = []byte(content)
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return files, ignoreFile, nil
}

func (d *Fetcher) fetchIgnoreListFromGitCommit(commit *object.Commit) (*bundle.IgnoreFile, error) {
	ignoreFile, err := commit.File(constant.IgnoreFileName)
	if err != nil {
		if err != object.ErrFileNotFound {
			return nil, err
		}
	}

	ignoreFileContent, err := ignoreFile.Contents()
	if err != nil {
		return nil, err
	}

	return d.Encoder.DecodeIgnoreFile([]byte(ignoreFileContent))
}
