package bpmiface

import (
	"github.com/4rchr4y/bpm/bundle"
	"github.com/4rchr4y/bpm/bundle/bundlefile"
	"github.com/4rchr4y/bpm/bundle/lockfile"
	"github.com/4rchr4y/bpm/bundleutil/encode"
)

type Encoder interface {
	DecodeIgnoreFile(content []byte) (*bundle.IgnoreFile, error)
	DecodeBundleFile(content []byte) (*bundlefile.Schema, error)
	DecodeLockFile(content []byte) (*lockfile.Schema, error)
	EncodeBundleFile(bundlefile *bundlefile.Schema) []byte
	EncodeLockFile(lockfile *lockfile.Schema) (result []byte)
	EncodeIgnoreFile(ignorefile *bundle.IgnoreFile) []byte
	Fileify(files map[string][]byte) (*encode.FileifyOutput, error)
}
