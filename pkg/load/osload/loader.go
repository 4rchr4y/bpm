package osload

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/pkg/bundle"
)

type osWrapper interface {
	Walk(root string, fn filepath.WalkFunc) error
	Open(name string) (*os.File, error)
}

type ioWrapper interface {
	ReadAll(reader io.Reader) ([]byte, error)
}

type bundleFileifier interface {
	Fileify(files map[string][]byte) (*bundle.Bundle, error)
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

func (loader *OsLoader) readBundleDir(dirPath string) (map[string][]byte, error) {
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

			relativePath := strings.TrimPrefix(path, filepath.Clean(dirPath)+"/")
			files[relativePath] = content

		}

		return nil
	}

	err := loader.osWrap.Walk(dirPath, walkFunc)
	if err != nil {
		return nil, fmt.Errorf("error walking the path %s: %v", dirPath, err)
	}

	return files, nil
}
