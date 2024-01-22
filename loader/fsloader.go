package loader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/4rchr4y/bpm/bundle"
	"github.com/4rchr4y/bpm/constant"
)

type fsLoaderOSWrapper interface {
	Walk(root string, fn filepath.WalkFunc) error
	Open(name string) (*os.File, error)
}

type fsLoaderIOWrapper interface {
	ReadAll(reader io.Reader) ([]byte, error)
}

// FsLoader is a files loader from a file system
type FsLoader struct {
	osWrap    fsLoaderOSWrapper
	ioWrap    fsLoaderIOWrapper
	fileifier bundleFileifier
}

func NewFsLoader(osWrap fsLoaderOSWrapper, ioWrap fsLoaderIOWrapper, fileifier bundleFileifier) *FsLoader {
	return &FsLoader{
		osWrap:    osWrap,
		ioWrap:    ioWrap,
		fileifier: fileifier,
	}
}

func (loader *FsLoader) LoadBundle(dirPath string) (*bundle.Bundle, error) {
	files, err := loader.readBundleDir(dirPath)
	if err != nil {
		return nil, err
	}

	if _, exist := files[constant.BundleFileName]; !exist {
		return nil, fmt.Errorf("'%s' is invalid bundle directory, can't find %s file", dirPath, constant.BundleFileName)
	}

	bundle, err := loader.fileifier.Fileify(files)
	if err != nil {
		return nil, err
	}

	return bundle, nil
}

func (loader *FsLoader) readBundleDir(dirPath string) (map[string][]byte, error) {
	files := make(map[string][]byte)
	// TODO: do file filtering
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error occurred while accessing a path %s: %v", path, err)
		}

		if !info.IsDir() {
			file, err := loader.osWrap.Open(path)
			if err != nil {
				return err
			}

			content, err := loader.ioWrap.ReadAll(file)
			if err != nil {
				return err
			}

			localPath := strings.Clone(path[len(dirPath)+1:])
			files[localPath] = content
		}

		return nil
	}

	err := loader.osWrap.Walk(dirPath, walkFunc)
	if err != nil {
		return nil, fmt.Errorf("error walking the path %s: %v", dirPath, err)
	}

	return files, nil
}
