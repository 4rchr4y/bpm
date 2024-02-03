package osload

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/bpm/pkg/bundleutil"
)

type osWrapper interface {
	Walk(root string, fn filepath.WalkFunc) error
	Open(name string) (*os.File, error)
}

type ioWrapper interface {
	ReadAll(reader io.Reader) ([]byte, error)
}

type bundleFileifier interface {
	Fileify(files map[string][]byte, options ...bundleutil.BundleOptFn) (*bundle.Bundle, error)
}

type OsLoader struct {
	osWrap    osWrapper
	ioWrap    ioWrapper
	fileifier bundleFileifier
}

func NewOsLoader(os osWrapper, io ioWrapper, fileifier bundleFileifier) *OsLoader {
	return &OsLoader{
		osWrap:    os,
		ioWrap:    io,
		fileifier: fileifier,
	}
}

func (loader *OsLoader) LoadBundle(dirPath string) (*bundle.Bundle, error) {
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

	bundle, err := loader.fileifier.Fileify(files, bundleutil.WithIgnoreList(ignoreList))
	if err != nil {
		return nil, err
	}

	return bundle, nil
}

func (loader *OsLoader) readBundleDir(abs string, ignoreList map[string]struct{}) (map[string][]byte, error) {
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

func (loader *OsLoader) readFileContent(path string) ([]byte, error) {
	file, err := loader.osWrap.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return loader.ioWrap.ReadAll(file)
}

func (loader *OsLoader) fetchIgnoreList(dir string) (map[string]struct{}, error) {
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
