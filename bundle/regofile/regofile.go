package regofile

import (
	"crypto/sha256"
	"strings"

	"github.com/4rchr4y/bpm/bundleutil"
	"github.com/open-policy-agent/opa/ast"
)

const ImportPathPrefix = "data."

type File struct {
	Path   string
	Raw    []byte
	Parsed *ast.Module
}

func (f *File) Package() string {
	return strings.TrimPrefix(f.Parsed.Package.Path.String(), ImportPathPrefix)
}

func (f *File) Sum() string {
	return bundleutil.ChecksumSHA256(sha256.New(), f.Parsed.String())
}
