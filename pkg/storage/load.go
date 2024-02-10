package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/bpm/pkg/bundleutil/bundlebuild"
)

type ErrNotExist struct{}

func (ErrNotExist) Error() string { return "bundle does not exist" }

func (s *Storage) Load(repo string, version *bundle.VersionExpr) (*bundle.Bundle, error) {
	return s.LoadFromAbs(s.MakeBundleSourcePath(repo, version.String()))
}

func (s *Storage) LoadFromAbs(dir string) (*bundle.Bundle, error) {
	ok, err := s.OSWrap.Exists(dir)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNotExist{}
	}

	absDirPath, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path for %s: %v", dir, err)
	}

	s.IO.PrintfInfo("loading from %s", absDirPath)

	ignoreList, err := s.readIgnoreFile(absDirPath)
	if err != nil {
		return nil, err
	}

	files, err := s.readBundleDir(absDirPath, ignoreList)
	if err != nil {
		return nil, err
	}

	bundle, err := s.Encoder.Fileify(files, bundlebuild.WithIgnoreList(ignoreList))
	if err != nil {
		return nil, err
	}

	return bundle, nil
}

func (fetcher *Storage) readBundleDir(abs string, ignoreFile *bundle.IgnoreFile) (map[string][]byte, error) {
	files := make(map[string][]byte)
	err := fetcher.OSWrap.Walk(abs, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error occurred while accessing a path %s: %v", path, err)
		}

		relativePath, err := filepath.Rel(abs, path)
		if err != nil {
			return fmt.Errorf("error getting relative path for %s from %s: %v", path, abs, err)
		}

		if ignoreFile != nil && ignoreFile.Lookup(relativePath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.IsDir() {
			content, err := fetcher.readFileContent(path)
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

func (fetcher *Storage) readFileContent(path string) ([]byte, error) {
	file, err := fetcher.OSWrap.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// TODO: use bufio buffer instead of ReadAll
	return fetcher.IOWrap.ReadAll(file)
}

func (fetcher *Storage) readIgnoreFile(dir string) (*bundle.IgnoreFile, error) {
	ignoreFilePath := filepath.Join(dir, constant.IgnoreFileName)
	file, err := fetcher.OSWrap.OpenFile(ignoreFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}
	defer file.Close()

	// TODO: use bufio buffer instead of ReadAll
	content, err := fetcher.IOWrap.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return fetcher.Encoder.DecodeIgnoreFile(content)
}
