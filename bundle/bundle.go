package bundle

import (
	"fmt"

	"github.com/4rchr4y/bpm/constant"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type VersionExpr struct {
	Version   string // semantic version if available, or pseudo version
	Timestamp string // commit timestamp
	Hash      string // commit hash
}

func NewVersionExpr(commit *object.Commit, tag string) *VersionExpr {
	timestamp := commit.Committer.When.UTC().Format("20060102150405")
	shortHash := commit.Hash.String()[:12]

	if tag == "" {
		tag = constant.BundlePseudoVersion
	}

	return &VersionExpr{
		Version:   tag,
		Timestamp: timestamp,
		Hash:      shortHash,
	}
}

func (v *VersionExpr) String() string {
	return fmt.Sprintf("%s+%s-%s", v.Version, v.Timestamp, v.Hash)
}

type Bundle struct {
	Version        *VersionExpr
	BundleFile     *BundleFile
	BundleLockFile *BundleLockFile
	BpmWorkFile    *BpmWorkFile
	RegoFiles      map[string]*RawRegoFile
}

func (b *Bundle) Name() string { return b.BundleFile.Package.Name }

func (b *Bundle) UpdateLock() bool {
	if len(b.RegoFiles) < 1 {
		// no rego files, then nothing to update
		return false
	}

	if b.BundleLockFile == nil {
		b.BundleLockFile = &BundleLockFile{
			Version: 1,
			Modules: make([]*ModuleDef, len(b.RegoFiles)),
		}
	}

	var i uint
	for path, file := range b.RegoFiles {
		b.BundleLockFile.Modules[i] = &ModuleDef{
			Name:     file.Package(),
			Source:   path,
			Checksum: file.Sum(),
			Dependencies: func() []string {
				result := make([]string, len(file.Parsed.Imports))
				for i, _import := range file.Parsed.Imports {
					result[i] = _import.Path.String()
				}

				return result
			}(),
		}

		i++
	}

	return true
}

func (b *Bundle) Validation() error {
	panic("not implemented")
}
