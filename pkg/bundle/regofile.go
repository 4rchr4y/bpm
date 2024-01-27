package bundle

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/open-policy-agent/opa/ast"
)

type RegoFile struct {
	Path   string
	Raw    []byte
	Parsed *ast.Module
}

func (rf *RegoFile) Package() string {
	return rf.Parsed.Package.Path.String()
}

func (rf *RegoFile) Sum() string {
	hash := md5.Sum([]byte(rf.Parsed.String()))
	return hex.EncodeToString(hash[:])
}
