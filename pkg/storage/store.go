package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/bpm/pkg/bundle/bundlefile"
	"github.com/4rchr4y/bpm/pkg/bundle/lockfile"
	"github.com/4rchr4y/bpm/pkg/bundle/regofile"
)

func (s *Storage) StoreMultiple(bundleset []*bundle.Bundle) error {
	for _, b := range bundleset {
		if err := s.Store(b); err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) Store(b *bundle.Bundle) error {
	dirPath := s.MakeBundleSourcePath(b.Repository(), b.Version.String())

	s.IO.PrintfDebug("saving to %s", dirPath)

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

func (s *Storage) processBundleLockFile(lockFile *lockfile.File, dir string) error {
	bytes := s.Encoder.EncodeLockFile(lockFile)
	path := filepath.Join(dir, lockFile.Filename())
	if err := s.OSWrap.WriteFile(path, bytes, 0644); err != nil {
		return err
	}

	return nil
}

func (s *Storage) processBundleFile(bundleFile *bundlefile.File, dir string) error {
	bytes := s.Encoder.EncodeBundleFile(bundleFile)
	path := filepath.Join(dir, bundleFile.Filename())
	if err := s.OSWrap.WriteFile(path, bytes, 0644); err != nil {
		return err
	}

	return nil
}

func (s *Storage) processRegoFiles(files map[string]*regofile.File, dir string) error {
	for filePath, file := range files {
		pathToSave := filepath.Join(dir, filePath)
		dirToSave := filepath.Dir(pathToSave)

		if _, err := os.Stat(dirToSave); os.IsNotExist(err) {
			if err := os.MkdirAll(dirToSave, 0755); err != nil {
				return fmt.Errorf("failed to create directory '%s': %v", dirToSave, err)
			}
		} else if err != nil {
			return fmt.Errorf("error checking directory '%s': %v", dirToSave, err)
		}

		if err := s.OSWrap.WriteFile(pathToSave, file.Raw, 0644); err != nil {
			return fmt.Errorf("failed to write file '%s': %v", pathToSave, err)
		}
	}

	return nil
}
