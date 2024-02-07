package bundleutil

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/bpm/pkg/bundle/bundlefile"
	"github.com/4rchr4y/bpm/pkg/bundle/lockfile"
	"github.com/4rchr4y/bpm/pkg/bundle/regofile"
	"github.com/4rchr4y/godevkit/syswrap/osiface"
)

type saverHclEncoder interface {
	EncodeBundleFile(bundlefile *bundlefile.File) []byte
	EncodeLockFile(lockfile *lockfile.File) []byte
}

type Saver struct {
	Dir    string // folder where all packages will be saved
	osWrap osiface.OSWrapper
	encode saverHclEncoder
}

func NewSaver(osWrap osiface.OSWrapper, encoder saverHclEncoder) *Saver {
	return &Saver{
		osWrap: osWrap,
		encode: encoder,
	}
}

func (s *Saver) SaveToDisk(bundles ...*bundle.Bundle) error {
	homeDir, err := s.osWrap.UserHomeDir()
	if err != nil {
		return err
	}

	for _, b := range bundles {
		if err := s.save(homeDir, b); err != nil {
			return fmt.Errorf("can't save bundle '%s': %v", b.Name(), err)
		}
	}

	return nil
}

func (s *Saver) save(homeDir string, b *bundle.Bundle) error {
	dirPath := filepath.Join(homeDir, constant.BPMDirName, b.Repository(), b.Version.String())

	if err := s.processRegoFiles(b.RegoFiles, dirPath); err != nil {
		return fmt.Errorf("error occurred rego files processing: %v", err)
	}

	if err := s.processBundleLockFile(b.LockFile, dirPath); err != nil {
		return fmt.Errorf("failed to encode %s file: %v", b.LockFile.Filename(), err)
	}

	if err := s.processBundleFile(b.BundleFile, dirPath); err != nil {
		return fmt.Errorf("failed to encode %s file: %v", b.BundleFile.Filename(), err)
	}

	return nil
}

func (s *Saver) processBundleLockFile(lockFile *lockfile.File, bundleVersionDir string) error {
	bytes := s.encode.EncodeLockFile(lockFile)
	path := filepath.Join(bundleVersionDir, lockFile.Filename())
	if err := s.osWrap.WriteFile(path, bytes, 0644); err != nil {
		return err
	}

	return nil
}

func (s *Saver) processBundleFile(bundleFile *bundlefile.File, bundleVersionDir string) error {
	bytes := s.encode.EncodeBundleFile(bundleFile)
	path := filepath.Join(bundleVersionDir, bundleFile.Filename())
	if err := s.osWrap.WriteFile(path, bytes, 0644); err != nil {
		return err
	}

	return nil
}

func (s *Saver) processRegoFiles(files map[string]*regofile.File, bundleVersionDir string) error {
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

		if err := s.osWrap.WriteFile(pathToSave, file.Raw, 0644); err != nil {
			return fmt.Errorf("failed to write file '%s': %v", pathToSave, err)
		}
	}

	return nil
}
