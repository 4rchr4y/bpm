package fetch

import (
	"context"
	"fmt"
	"strings"

	"github.com/4rchr4y/bpm/bundle"
	"github.com/4rchr4y/bpm/bundle/bundlefile"
	"github.com/4rchr4y/bpm/bundle/lockfile"
	"github.com/4rchr4y/bpm/bundleutil"
	"github.com/4rchr4y/bpm/bundleutil/encode"
	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/iostream/iostreamiface"
	"github.com/4rchr4y/bpm/regoutil"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hashicorp/go-version"
)

type githubFetcherEncoder interface {
	DecodeBundleFile(content []byte) (*bundlefile.Schema, error)
	DecodeIgnoreFile(content []byte) (*bundle.IgnoreFile, error)
	DecodeLockFile(content []byte) (*lockfile.Schema, error)
	Fileify(files map[string][]byte) (*encode.FileifyOutput, error)
}

type githubFetcherClient interface {
	CloneWithContext(ctx context.Context, opts *git.CloneOptions) (*git.Repository, error)
}

type GithubFetcher struct {
	IO      iostreamiface.IO
	Client  githubFetcherClient
	Encoder githubFetcherEncoder
}

func (gh *GithubFetcher) Download(ctx context.Context, source string, tag *bundle.VersionSpec) (*bundle.Bundle, error) {
	gh.IO.PrintfInfo("downloading %s", bundleutil.FormatSourceWithVersion(source, tag.String()))

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

	filesOutput, err := gh.getFilesFromCommit(commit)
	if err != nil {
		return nil, err
	}

	fileifyOutput, err := gh.Encoder.Fileify(filesOutput.FileSet)
	if err != nil {
		return nil, err
	}

	return &bundle.Bundle{
		Version:    v,
		Source:     source,
		BundleFile: bundlefile.PrepareSchema(filesOutput.BundleFile),
		LockFile:   lockfile.PrepareSchema(filesOutput.LockFile),
		RegoFiles:  fileifyOutput.RegoFiles,
		IgnoreFile: filesOutput.IgnoreFile,
		OtherFiles: fileifyOutput.OtherFiles,
	}, nil
}

type getFilesOutput struct {
	FileSet    map[string][]byte
	BundleFile *bundlefile.Schema
	LockFile   *lockfile.Schema
	IgnoreFile *bundle.IgnoreFile
}

func (gh *GithubFetcher) getFilesFromCommit(commit *object.Commit) (output *getFilesOutput, err error) {
	output = new(getFilesOutput)

	output.IgnoreFile, err = gh.readIgnoreFileFromGitCommit(commit)
	if err != nil {
		return nil, err
	}

	output.BundleFile, err = gh.readBundleFileFromGitCommit(commit)
	if err != nil {
		return nil, err
	}

	output.LockFile, err = gh.readLockFileFromGitCommit(commit)
	if err != nil {
		return nil, err
	}

	if err := regoutil.PrepareDocumentParser(output.BundleFile); err != nil {
		return nil, err
	}

	filesIter, err := commit.Files()
	if err != nil {
		return nil, err
	}

	output.FileSet = make(map[string][]byte)
	err = filesIter.ForEach(func(f *object.File) error {
		if output.IgnoreFile.Some(f.Name) {
			return nil
		}

		content, err := f.Contents()
		if err != nil {
			return err
		}

		output.FileSet[f.Name] = []byte(content)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (gh *GithubFetcher) readIgnoreFileFromGitCommit(commit *object.Commit) (*bundle.IgnoreFile, error) {
	content, err := readFileFromCommit(commit, constant.IgnoreFileName)
	if err != nil {
		if err != object.ErrFileNotFound {
			return nil, err
		}

		return nil, nil
	}

	return gh.Encoder.DecodeIgnoreFile([]byte(content))
}

func (gh *GithubFetcher) readLockFileFromGitCommit(commit *object.Commit) (*lockfile.Schema, error) {
	content, err := readFileFromCommit(commit, constant.LockFileName)
	if err != nil {
		if err != object.ErrFileNotFound {
			return nil, err
		}

		return nil, nil
	}

	return gh.Encoder.DecodeLockFile([]byte(content))
}

func (gh *GithubFetcher) readBundleFileFromGitCommit(commit *object.Commit) (*bundlefile.Schema, error) {
	content, err := readFileFromCommit(commit, constant.BundleFileName)
	if err != nil {
		return nil, err
	}

	return gh.Encoder.DecodeBundleFile([]byte(content))
}

func (gh *GithubFetcher) fetchCommitByTag(repo *git.Repository, v *bundle.VersionSpec) (*object.Commit, *bundle.VersionSpec, error) {
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

func (gh *GithubFetcher) getPseudoVersionCommit(repo *git.Repository, v *bundle.VersionSpec) (*object.Commit, error) {
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

func (gh *GithubFetcher) getLatestVersionCommit(repo *git.Repository) (*object.Commit, *bundle.VersionSpec, error) {
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

		return commit, bundle.NewVersionSpecFromCommit(commit, v), nil
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, nil, err
	}

	return commit, bundle.NewVersionSpecFromCommit(commit, v), nil
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

func readFileFromCommit(commit *object.Commit, fileName string) ([]byte, error) {
	f, err := commit.File(fileName)
	if err != nil {
		return nil, err
	}

	content, err := f.Contents()
	if err != nil {
		return nil, err
	}

	return []byte(content), nil
}
