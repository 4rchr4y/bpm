package gitload

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hashicorp/go-version"

	"github.com/4rchr4y/bpm/pkg/bundle"
)

type bundleFileifier interface {
	Fileify(files map[string][]byte) (*bundle.Bundle, error)
}

type gitClient interface {
	CloneWithContext(ctx context.Context, opts *git.CloneOptions) (*git.Repository, error)
}

type GitLoader struct {
	fileifier bundleFileifier
	gitCli    gitClient
}

func NewGitLoader(gitClient gitClient, bparser bundleFileifier) *GitLoader {
	return &GitLoader{
		gitCli:    gitClient,
		fileifier: bparser,
	}
}

func (loader *GitLoader) DownloadBundle(ctx context.Context, url string, tag string) (*bundle.Bundle, error) {
	cloneInput := &git.CloneOptions{
		URL: fmt.Sprintf("https://%s.git", url),
	}

	repo, err := loader.gitCli.CloneWithContext(ctx, cloneInput)
	if err != nil {
		return nil, err
	}

	commit, v, err := fetchCommitByTag(repo, tag)
	if err != nil {
		return nil, err
	}

	files, err := getFilesFromCommit(commit)
	if err != nil {
		return nil, err
	}

	b, err := loader.fileifier.Fileify(files)
	if err != nil {
		return nil, err
	}

	b.Version = v
	b.BundleFile.Package.Repository = url

	return b, nil
}

func fetchCommitByTag(repo *git.Repository, tag string) (*object.Commit, *bundle.VersionExpr, error) {
	v, err := bundle.ParseVersionExpr(tag)
	if err != nil {
		switch err {
		case bundle.ErrVersionInvalidFormat{}:
			return nil, nil, err

		case bundle.ErrEmptyVersion{}:
			return getLatestVersionCommit(repo)

		default:
			return nil, nil, err
		}
	}

	if v.IsPseudo() {
		commit, err := getPseudoVersionCommit(repo, v)
		if err != nil {
			return nil, nil, err
		}

		return commit, v, nil
	}

	commit, err := getCurrentVersionCommit(repo, tag)
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
		if strings.HasPrefix(c.Hash.String(), v.Hash) && v.Timestamp == c.Committer.When.UTC().Format(bundle.VersionDateFormat) {
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

	v, ref := findLatestVersion(tags)
	if v == nil || ref == nil {
		commit, err := getLatestCommit(repo)
		if err != nil {
			return nil, nil, err
		}

		return commit, bundle.NewVersionExpr(commit, v), nil
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, nil, err
	}

	return commit, bundle.NewVersionExpr(commit, v), nil
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

func findLatestVersion(tags map[*version.Version]*plumbing.Reference) (v *version.Version, ref *plumbing.Reference) {
	for version, reference := range tags {
		if v == nil || version.GreaterThan(v) {
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

func getFilesFromCommit(commit *object.Commit) (map[string][]byte, error) {
	filesIter, err := commit.Files()
	if err != nil {
		return nil, err
	}

	files := make(map[string][]byte)
	err = filesIter.ForEach(func(f *object.File) error {
		content, err := f.Contents()
		if err != nil {
			return err
		}
		files[f.Name] = []byte(content)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
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
