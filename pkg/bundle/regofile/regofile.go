package regofile

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/open-policy-agent/opa/ast"
)

type File struct {
	Path   string
	Raw    []byte
	Parsed *ast.Module
}

func (f *File) Package() string {
	return f.Parsed.Package.Path.String()
}

func (f *File) Sum() string {
	hash := md5.Sum([]byte(f.Parsed.String()))
	return hex.EncodeToString(hash[:])
}
