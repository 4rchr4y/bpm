package bundleutil

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hashicorp/go-version"
)

type loaderOsWrapper interface {
	Walk(root string, fn filepath.WalkFunc) error
	Open(name string) (*os.File, error)
}

type loaderIoWrapper interface {
	ReadAll(reader io.Reader) ([]byte, error)
}

type loaderFileifier interface {
	Fileify(files map[string][]byte, options ...BundleOptFn) (*bundle.Bundle, error)
}
type loaderGitFacade interface {
	CloneWithContext(ctx context.Context, opts *git.CloneOptions) (*git.Repository, error)
}

type Loader struct {
	osWrap    loaderOsWrapper
	ioWrap    loaderIoWrapper
	fileifier loaderFileifier
	gitfacade loaderGitFacade
}

func NewLoader(os loaderOsWrapper, io loaderIoWrapper, fileifier loaderFileifier, gitfacade loaderGitFacade) *Loader {
	return &Loader{
		osWrap:    os,
		ioWrap:    io,
		fileifier: fileifier,
		gitfacade: gitfacade,
	}
}

// -----------------------------------------------------------------------
// Loader of bundle files from the file system

func (loader *Loader) LoadBundle(dirPath string) (*bundle.Bundle, error) {
	absDirPath, err := filepath.Abs(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path for %s: %v", dirPath, err)
	}

	ignoreList, err := loader.fetchIgnoreList(absDirPath)
	if err != nil {
		return nil, err
	}

	files, err := loader.readBundleDir(absDirPath, ignoreList)
	if err != nil {
		return nil, err
	}

	bundle, err := loader.fileifier.Fileify(files, WithIgnoreList(ignoreList))
	if err != nil {
		return nil, err
	}

	return bundle, nil
}

func (loader *Loader) readBundleDir(abs string, ignoreList map[string]struct{}) (map[string][]byte, error) {
	files := make(map[string][]byte)
	err := loader.osWrap.Walk(abs, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error occurred while accessing a path %s: %v", path, err)
		}

		relativePath, err := filepath.Rel(abs, path)
		if err != nil {
			return fmt.Errorf("error getting relative path for %s from %s: %v", path, abs, err)
		}

		if shouldIgnore(ignoreList, relativePath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.IsDir() {
			content, err := loader.readFileContent(path)
			if err != nil {
				return err
			}
			files[relativePath] = content
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking the path %s: %v", abs, err)
	}

	return files, nil
}

func (loader *Loader) readFileContent(path string) ([]byte, error) {
	file, err := loader.osWrap.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return loader.ioWrap.ReadAll(file)
}

func (loader *Loader) fetchIgnoreList(dir string) (map[string]struct{}, error) {
	ignoreFilePath := filepath.Join(dir, constant.IgnoreFileName)
	file, err := loader.osWrap.Open(ignoreFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]struct{}), nil
		}
		return nil, err
	}
	defer file.Close()

	content, err := loader.ioWrap.ReadAll(file)
	if err != nil {
		return nil, err
	}

	result := make(map[string]struct{})
	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		result[line] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading '%s' input: %v", constant.IgnoreFileName, err)
	}

	return result, nil
}

func shouldIgnore(ignoreList map[string]struct{}, path string) bool {
	if path == "" || len(ignoreList) == 0 {
		return false
	}

	dir := filepath.Dir(path)
	if dir == "." {
		return false
	}

	topLevelDir := strings.Split(dir, string(filepath.Separator))[0]
	_, found := ignoreList[topLevelDir]
	return found
}

// -----------------------------------------------------------------------
// Bundle downloader from a remote git server

func (loader *Loader) DownloadBundle(ctx context.Context, url string, tag string) (*bundle.Bundle, error) {
	cloneInput := &git.CloneOptions{
		URL: fmt.Sprintf("https://%s.git", url),
	}

	repo, err := loader.gitfacade.CloneWithContext(ctx, cloneInput)
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

	ignoreFile, err := commit.File("constant.IgnoreFileName")
	if err != nil {
		if err != object.ErrFileNotFound {
			return nil, err
		}
	}

	fmt.Println(ignoreFile.Name)

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
