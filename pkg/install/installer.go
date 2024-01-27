package install

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/4rchr4y/bpm/bundle"
	"github.com/4rchr4y/bpm/bundle/bundlefile"
	"github.com/4rchr4y/bpm/bundle/lockfile"
	"github.com/4rchr4y/bpm/constant"
)

type osWrapper interface {
	Mkdir(name string, perm fs.FileMode) error
	MkdirAll(path string, perm fs.FileMode) error
	Create(name string) (*os.File, error)
	UserHomeDir() (string, error)
	WriteFile(name string, data []byte, perm fs.FileMode) error
}

type hclEncoder interface {
	EncodeBundleFile(bundlefile *bundlefile.File) []byte
	EncodeLockFile(lockfile *lockfile.File) []byte
}

type BundleInstaller struct {
	osWrap osWrapper
	encode hclEncoder
}

func NewBundleInstaller(osWrap osWrapper, encoder hclEncoder) *BundleInstaller {
	return &BundleInstaller{
		osWrap: osWrap,
		encode: encoder,
	}
}

func (cmd *BundleInstaller) Install(b *bundle.Bundle) error {
	homeDir, err := cmd.osWrap.UserHomeDir()
	if err != nil {
		return err
	}

	dirPath := filepath.Join(homeDir, constant.BPMDirName, b.BundleFile.Package.Name, b.Version.String())

	if err := cmd.processRegoFiles(b.RegoFiles, dirPath); err != nil {
		return fmt.Errorf("error occurred rego files processing: %v", err)
	}

	if err := cmd.processBundleLockFile(b.BundleLockFile, dirPath); err != nil {
		return fmt.Errorf("failed to encode %s file: %v", b.BundleLockFile.FileName(), err)
	}

	if err := cmd.processBundleFile(b.BundleFile, dirPath); err != nil {
		return fmt.Errorf("failed to encode %s file: %v", b.BundleFile.FileName(), err)
	}

	return nil
}

func (cmd *BundleInstaller) processBundleLockFile(bundleLockFile *lockfile.File, bundleVersionDir string) error {
	bytes := cmd.encode.EncodeLockFile(bundleLockFile)
	path := filepath.Join(bundleVersionDir, constant.LockFileName)
	if err := cmd.osWrap.WriteFile(path, bytes, 0644); err != nil {
		return err
	}

	return nil
}

func (cmd *BundleInstaller) processBundleFile(bundleFile *bundlefile.File, bundleVersionDir string) error {
	bytes := cmd.encode.EncodeBundleFile(bundleFile)
	path := filepath.Join(bundleVersionDir, constant.BundleFileName)
	if err := cmd.osWrap.WriteFile(path, bytes, 0644); err != nil {
		return err
	}

	return nil
}

func (cmd *BundleInstaller) processRegoFiles(files map[string]*bundle.RawRegoFile, bundleVersionDir string) error {
	for filePath, file := range files {
		pathToSave := filepath.Join(bundleVersionDir, filePath)
		dirToSave := filepath.Dir(pathToSave)

		if _, err := os.Stat(dirToSave); os.IsNotExist(err) {
			if err := os.MkdirAll(dirToSave, 0755); err != nil {
				return fmt.Errorf("failed to create directory '%s': %v", dirToSave, err)
			}
		} else if err != nil {
			return fmt.Errorf("error checking directory '%s': %v", dirToSave, err)
		}

		if err := cmd.osWrap.WriteFile(pathToSave, file.Raw, 0644); err != nil {
			return fmt.Errorf("failed to write file '%s': %v", pathToSave, err)
		}
	}

	return nil
}
