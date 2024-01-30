package gitload

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/hashicorp/go-version"

	gitcli "github.com/4rchr4y/bpm/internal/git"
	"github.com/4rchr4y/bpm/pkg/bundle"
)

type bundleFileifier interface {
	Fileify(files map[string][]byte) (*bundle.Bundle, error)
}

type gitClient interface {
	PlainClone(input *gitcli.PlainCloneInput) (*git.Repository, error)
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

func (loader *GitLoader) DownloadBundle(url string, versionStr string) (*bundle.Bundle, error) {
	repoURL := fmt.Sprintf("https://%s.git", url)
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{URL: repoURL})
	if err != nil {
		return nil, err
	}

	var commit *object.Commit

	v, err := bundle.ParseVersionExpr(versionStr)
	if err != nil {
		switch err {
		case bundle.ErrVersionInvalidFormat{}:
			return nil, err

		case bundle.ErrEmptyVersion{}:
			commit, v, err = getLatestVersionCommit(repo)
			if err != nil {
				return nil, err
			}

		default:
			return nil, err
		}
	}

	switch {
	case v.IsPseudo():
		commit, err = getPseudoVersionCommit(repo, v)
	}

	// if versionStr != "" {
	// 	ref, err := getRef(repo, versionStr)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	commit, err = repo.CommitObject(ref)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// } else {
	// 	commit, err = getLatestCommit(repo)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	files, err := getFilesFromCommit(commit)
	if err != nil {
		return nil, err
	}

	b, err := loader.fileifier.Fileify(files)
	if err != nil {
		return nil, err
	}

	b.Version = v // bundle.NewVersionExpr(commit, versionStr)
	b.BundleFile.Package.Repository = url

	return b, nil
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

func findLatestVersion(tags map[*version.Version]*plumbing.Reference) (*version.Version, *plumbing.Reference) {
	var (
		latestVersion   *version.Version
		latestReference *plumbing.Reference
	)

	for v, ref := range tags {
		if latestVersion == nil || v.GreaterThan(latestVersion) {
			latestVersion = v
			latestReference = ref
		}
	}

	return latestVersion, latestReference
}

func getLatestCommit(repo *git.Repository) (*object.Commit, error) {
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}
	return repo.CommitObject(ref.Hash())
}

// func getRef(repo *git.Repository, versionRaw string) (plumbing.Hash, error) {
// 	version, err := bundle.ParseVersionExpr(versionRaw)
// 	if err != nil {
// 		// return nil, err
// 	}

// 	// fmt.Println(version.Tag, version.Hash)

// 	if version.IsPseudo() {
// 		fmt.Println(versionRaw)

// 		commit, err := findCommitByShortHash(repo, version)
// 		if err != nil {
// 			return plumbing.Hash{}, fmt.Errorf("version by hash is not found: %s", versionRaw)
// 		}

// 		return commit.Hash, nil
// 	}

// 	ref, err := repo.Tag(version.Version)
// 	if err != nil {
// 		return plumbing.Hash{}, fmt.Errorf("version by tag is not found: %s", versionRaw)
// 	}

// 	return ref.Hash(), nil
// }

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

// func findCommitByShortHash(repo *git.Repository, version *bundle.VersionExpr) (*object.Commit, error) {
// 	if len(version.Hash) != bundle.VersionShortHashLen {
// 		return nil, fmt.Errorf("short hash must be %d characters long", bundle.VersionShortHashLen)
// 	}

// 	iter, err := repo.CommitObjects()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer iter.Close()

// 	var foundCommit *object.Commit
// 	err = iter.ForEach(func(c *object.Commit) error {
// 		if strings.HasPrefix(c.Hash.String(), version.Hash) && version.Timestamp == c.Committer.When.UTC().Format(bundle.VersionDateFormat) {
// 			foundCommit = c
// 			return nil
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	if foundCommit == nil {
// 		return nil, fmt.Errorf("commit not found with hash %s", version.Hash)
// 	}

// 	return foundCommit, nil
// }

func findTag(repo *git.Repository, tag string) (*plumbing.Reference, error) {
	tags, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	var foundTag *plumbing.Reference
	err = tags.ForEach(func(t *plumbing.Reference) error {
		if t.Name().Short() == tag {
			foundTag = t
			return nil
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if foundTag == nil {
		return nil, fmt.Errorf("version '%s' is not found", tag)
	}

	return foundTag, nil
}
