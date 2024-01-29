package regofile

import (
	"crypto/sha256"

	"github.com/4rchr4y/bpm/pkg/util"
	"github.com/open-policy-agent/opa/ast"
)

type File struct {
	Path   string
	Raw    []byte
	Parsed *ast.Module
}

func (f *File) Package() string { return f.Parsed.Package.Path.String() }
func (f *File) Sum() string     { return util.ChecksumSHA256(sha256.New(), f.Parsed.String()) }
