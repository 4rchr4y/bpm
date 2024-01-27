package bundle

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/open-policy-agent/opa/ast"
)

type RawRegoFile struct {
	Path   string
	Raw    []byte
	Parsed *ast.Module
}

func (rrf *RawRegoFile) Package() string {
	return rrf.Parsed.Package.Path.String()
}

func (rrf *RawRegoFile) Sum() string {
	hash := md5.Sum([]byte(rrf.Parsed.String()))
	return hex.EncodeToString(hash[:])
}
