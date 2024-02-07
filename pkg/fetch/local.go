package fetch

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/bpm/pkg/bundleutil"
	"github.com/4rchr4y/bpm/pkg/fileutil"
)

func (fetcher *Fetcher) FetchLocal(dirPath string) (*bundle.Bundle, error) {
	absDirPath, err := filepath.Abs(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path for %s: %v", dirPath, err)
	}

	fetcher.IO.PrintfInfo("loading bundle from %s", absDirPath)

	ignoreList, err := fetcher.fetchIgnoreListFromLocalFile(absDirPath)
	if err != nil {
		return nil, err
	}

	files, err := fetcher.readLocalBundleDir(absDirPath, ignoreList)
	if err != nil {
		return nil, err
	}

	bundle, err := fetcher.Fileifier.Fileify(files, bundleutil.WithIgnoreList(ignoreList))
	if err != nil {
		return nil, err
	}

	return bundle, nil
}

func (fetcher *Fetcher) readLocalBundleDir(abs string, ignoreList map[string]struct{}) (map[string][]byte, error) {
	files := make(map[string][]byte)
	err := fetcher.OSWrap.Walk(abs, func(path string, info os.FileInfo, err error) error {
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
			content, err := fetcher.readLocalFileContent(path)
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

func (fetcher *Fetcher) readLocalFileContent(path string) ([]byte, error) {
	file, err := fetcher.OSWrap.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return fetcher.IOWrap.ReadAll(file)
}

func (fetcher *Fetcher) fetchIgnoreListFromLocalFile(dir string) (map[string]struct{}, error) {
	ignoreFilePath := filepath.Join(dir, constant.IgnoreFileName)
	file, err := fetcher.OSWrap.OpenFile(ignoreFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]struct{}), nil
		}
		return nil, err
	}
	defer file.Close()

	content, err := fetcher.IOWrap.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return fileutil.ReadLinesToMap(content)
}
