package install

import (
	"io/fs"
	"os"

	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/bpm/pkg/bundle/bundlefile"
	"github.com/4rchr4y/bpm/pkg/bundle/lockfile"
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

type gitLoader interface {
	DownloadBundle(url string, tag string) (*bundle.Bundle, error)
}

// type BundleInstaller struct {
// 	git    gitLoader
// 	osWrap osWrapper
// 	encode hclEncoder
// }

// func NewBundleInstaller(osWrap osWrapper, encoder hclEncoder) *BundleInstaller {
// 	return &BundleInstaller{
// 		osWrap: osWrap,
// 		encode: encoder,
// 	}
// }

// func (bi *BundleInstaller) SaveToDisk(b *bundle.Bundle) error {
// 	homeDir, err := bi.osWrap.UserHomeDir()
// 	if err != nil {
// 		return err
// 	}

// 	dirPath := filepath.Join(homeDir, constant.BPMDirName, b.Repository(), b.Version.String())

// 	if err := bi.processRegoFiles(b.RegoFiles, dirPath); err != nil {
// 		return fmt.Errorf("error occurred rego files processing: %v", err)
// 	}

// 	if err := bi.processBundleLockFile(b.LockFile, dirPath); err != nil {
// 		return fmt.Errorf("failed to encode %s file: %v", b.LockFile.Filename(), err)
// 	}

// 	if err := bi.processBundleFile(b.BundleFile, dirPath); err != nil {
// 		return fmt.Errorf("failed to encode %s file: %v", b.BundleFile.Filename(), err)
// 	}

// 	return nil
// }

// func (bi *BundleInstaller) processBundleLockFile(lockFile *lockfile.File, bundleVersionDir string) error {
// 	bytes := bi.encode.EncodeLockFile(lockFile)
// 	path := filepath.Join(bundleVersionDir, lockFile.Filename())
// 	if err := bi.osWrap.WriteFile(path, bytes, 0644); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (bi *BundleInstaller) processBundleFile(bundleFile *bundlefile.File, bundleVersionDir string) error {
// 	bytes := bi.encode.EncodeBundleFile(bundleFile)
// 	path := filepath.Join(bundleVersionDir, bundleFile.Filename())
// 	if err := bi.osWrap.WriteFile(path, bytes, 0644); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (bi *BundleInstaller) processRegoFiles(files map[string]*regofile.File, bundleVersionDir string) error {
// 	for filePath, file := range files {
// 		pathToSave := filepath.Join(bundleVersionDir, filePath)
// 		dirToSave := filepath.Dir(pathToSave)

// 		if _, err := os.Stat(dirToSave); os.IsNotExist(err) {
// 			if err := os.MkdirAll(dirToSave, 0755); err != nil {
// 				return fmt.Errorf("failed to create directory '%s': %v", dirToSave, err)
// 			}
// 		} else if err != nil {
// 			return fmt.Errorf("error checking directory '%s': %v", dirToSave, err)
// 		}

// 		if err := bi.osWrap.WriteFile(pathToSave, file.Raw, 0644); err != nil {
// 			return fmt.Errorf("failed to write file '%s': %v", pathToSave, err)
// 		}
// 	}

// 	return nil
// }
