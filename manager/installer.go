package manager

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/4rchr4y/bpm/bundle"
	"github.com/4rchr4y/bpm/constant"
)

type installerOSWrapper interface {
	Create(name string) (*os.File, error)
	WriteFile(name string, data []byte, perm fs.FileMode) error
}

type installerTOMLEncoder interface {
	Encode(value interface{}) ([]byte, error)
}

type BundleInstaller struct {
	osWrap  installerOSWrapper
	encoder installerTOMLEncoder
}

func NewBundleInstaller(osWrap installerOSWrapper, encoder installerTOMLEncoder) *BundleInstaller {
	return &BundleInstaller{
		osWrap:  osWrap,
		encoder: encoder,
	}
}

type BundleInstallInput struct {
	HomeDir string
	Bundle  *bundle.Bundle
}

func (cmd *BundleInstaller) Install(input *BundleInstallInput) error {
	dirPath := fmt.Sprintf("%s/%s/%s/%s", input.HomeDir, constant.BPMDirName, input.Bundle.BundleFile.Package.Name, input.Bundle.Version.String())

	if err := cmd.processRegoFiles(input.Bundle.RegoFiles, dirPath); err != nil {
		return fmt.Errorf("error occurred rego files processing: %v", err)
	}

	if err := cmd.processBundleLockFile(input.Bundle.BundleLockFile, dirPath); err != nil {
		return fmt.Errorf("failed to encode %s file: %v", input.Bundle.BundleLockFile.FileName(), err)
	}

	if err := cmd.processBundleFile(input.Bundle.BundleFile, dirPath); err != nil {
		return fmt.Errorf("failed to encode %s file: %v", input.Bundle.BundleFile.FileName(), err)
	}

	return nil
}

func (cmd *BundleInstaller) processBundleLockFile(bundleLockFile *bundle.BundleLockFile, bundleVersionDir string) error {
	// file, err := cmd.osWrap.Create(fmt.Sprintf("%s/%s", bundleVersionDir, constant.LockFileName))
	// if err != nil {
	// 	return err
	// }

	// if err := cmd.encoder.Encode(file, bundleLockFile); err != nil {
	// 	return err
	// }

	bytes, err := cmd.encoder.Encode(bundleLockFile)
	if err != nil {
		return err
	}

	path := filepath.Join(bundleVersionDir, constant.LockFileName)
	if err := cmd.osWrap.WriteFile(path, bytes, 0644); err != nil {
		return err
	}

	return nil
}

func (cmd *BundleInstaller) processBundleFile(bundleFile *bundle.BundleFile, bundleVersionDir string) error {
	// file, err := cmd.osWrap.Create(fmt.Sprintf("%s/%s", bundleVersionDir, constant.BundleFileName))
	// if err != nil {
	// 	return err
	// }

	// if err := cmd.encoder.Encode(file, bundleFile); err != nil {
	// 	return err
	// }

	bytes, err := cmd.encoder.Encode(bundleFile)
	if err != nil {
		return err
	}

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
