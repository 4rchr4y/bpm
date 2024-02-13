package storage

import (
	"fmt"
	"path/filepath"

	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/bpm/pkg/bundle/bundlefile"
	"github.com/4rchr4y/bpm/pkg/bundle/lockfile"
	"github.com/4rchr4y/godevkit/v3/syswrap/ioiface"
	"github.com/4rchr4y/godevkit/v3/syswrap/osiface"
)

type storageFetcher interface {
	FetchLocal(dirPath string) (*bundle.Bundle, error)
}

type storageHCLEncoder interface {
	EncodeBundleFile(bundlefile *bundlefile.File) []byte
	EncodeLockFile(lockfile *lockfile.File) []byte
	DecodeIgnoreFile(content []byte) (*bundle.IgnoreFile, error)
	Fileify(files map[string][]byte) (*bundle.BundleRaw, error)
}

type Storage struct {
	Dir string

	IO     core.IO
	OSWrap osiface.OSWrapper
	IOWrap ioiface.IOWrapper

	Encoder storageHCLEncoder
}

func (s *Storage) Some(repo string, version string) bool {
	ok, _ := s.OSWrap.Exists(s.MakeBundleSourcePath(repo, version))
	return ok
}

func (s *Storage) MakeBundleSourcePath(repo string, version string) string {
	return filepath.Join(s.Dir, fmt.Sprintf("%s@%s", repo, version))
}
