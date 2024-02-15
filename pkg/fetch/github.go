package fetch

import (
	"context"
	"fmt"
	"strings"

	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hashicorp/go-version"
)

type githubFetcherEncoder interface {
	DecodeIgnoreFile(content []byte) (*bundle.IgnoreFile, error)
	Fileify(files map[string][]byte) (*bundle.BundleRaw, error)
}

type githubFetcherClient interface {
	CloneWithContext(ctx context.Context, opts *git.CloneOptions) (*git.Repository, error)
}

type GithubFetcher struct {
	IO      core.IO
	Client  githubFetcherClient
	Encoder githubFetcherEncoder
}

func (gh *GithubFetcher) Download(ctx context.Context, source string, tag *bundle.VersionExpr) (*bundle.Bundle, error) {
	gh.IO.PrintfInfo("downloading %s@%s", source, tag.String())

	options := &git.CloneOptions{
		URL: fmt.Sprintf("https://%s.git", source),
	}

	repo, err := gh.Client.CloneWithContext(ctx, options)
	if err != nil {
		return nil, err
	}

	commit, v, err := gh.fetchCommitByTag(repo, tag)
	if err != nil {
		return nil, err
	}

	files, ignoreFile, err := gh.getFilesFromCommit(commit)
	if err != nil {
		return nil, err
	}

	bundleRaw, err := gh.Encoder.Fileify(files)
	if err != nil {
		return nil, err
	}

	return bundleRaw.ToBundle(v, ignoreFile)
}

func (gh *GithubFetcher) getFilesFromCommit(commit *object.Commit) (map[string][]byte, *bundle.IgnoreFile, error) {
	filesIter, err := commit.Files()
	if err != nil {
		return nil, nil, err
	}

	ignoreFile, err := gh.fetchIgnoreListFromGitCommit(commit)
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

func (gh *GithubFetcher) fetchIgnoreListFromGitCommit(commit *object.Commit) (*bundle.IgnoreFile, error) {
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

	return gh.Encoder.DecodeIgnoreFile([]byte(ignoreFileContent))
}

func (gh *GithubFetcher) fetchCommitByTag(repo *git.Repository, v *bundle.VersionExpr) (*object.Commit, *bundle.VersionExpr, error) {
	if v == nil {
		return gh.getLatestVersionCommit(repo)
	}

	if v.IsPseudo() {
		commit, err := gh.getPseudoVersionCommit(repo, v)
		if err != nil {
			return nil, nil, err
		}

		return commit, v, nil
	}

	commit, err := gh.getCurrentVersionCommit(repo, v.SemTag.Original())
	if err != nil {
		return nil, nil, err
	}

	return commit, v, nil

}

func (gh *GithubFetcher) getPseudoVersionCommit(repo *git.Repository, v *bundle.VersionExpr) (*object.Commit, error) {
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

func (gh *GithubFetcher) getLatestVersionCommit(repo *git.Repository) (*object.Commit, *bundle.VersionExpr, error) {
	tags, err := gh.collectTagList(repo)
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

func (gh *GithubFetcher) collectTagList(repo *git.Repository) (map[*version.Version]*plumbing.Reference, error) {
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

func (gh *GithubFetcher) getCurrentVersionCommit(repo *git.Repository, tag string) (*object.Commit, error) {
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

func mustFindLatestVersion(tags map[*version.Version]*plumbing.Reference) (v *version.Version, ref *plumbing.Reference) {
	v = bundle.PseudoSemTag

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
