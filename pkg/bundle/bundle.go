package bundle

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"sort"
	"strings"

	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/pkg/bundle/bundlefile"
	"github.com/4rchr4y/bpm/pkg/bundle/lockfile"
	"github.com/4rchr4y/bpm/pkg/bundle/regofile"
	"github.com/emirpasic/gods/v2/maps/hashbidimap"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	dateFormat      = "20060102150405"
	shortHashLength = 12
)

type VersionExpr struct {
	Version   string // semantic version if available, or pseudo version
	Timestamp string // commit timestamp
	Hash      string // commit hash
}

func NewVersionExpr(commit *object.Commit, tag string) *VersionExpr {
	if strings.TrimSpace(tag) == "" {
		tag = constant.BundlePseudoVersion
	}

	return &VersionExpr{
		Version:   tag,
		Timestamp: commit.Committer.When.UTC().Format(dateFormat),
		Hash:      commit.Hash.String()[:shortHashLength],
	}
}

func (v *VersionExpr) String() string {
	if v.Version != constant.BundlePseudoVersion {
		return v.Version
	}

	return fmt.Sprintf("%s+%s-%s", v.Version, v.Timestamp, v.Hash)
}

type AuthorExpr struct {
	Username string // value of git 'config --get user.username'
	Email    string // value of git 'config --get user.email'
}

func (author *AuthorExpr) String() string {
	return fmt.Sprintf("%s %s", author.Username, author.Email)
}

type Bundle struct {
	Version     *VersionExpr
	BundleFile  *bundlefile.File
	LockFile    *lockfile.File
	IgnoreFiles map[string]struct{}
	RegoFiles   map[string]*regofile.File
	OtherFiles  map[string][]byte

	requireCache *hashbidimap.Map[string, string]
}

func (b *Bundle) Name() string       { return b.BundleFile.Package.Name }
func (b *Bundle) Repository() string { return b.BundleFile.Package.Repository }

func (b *Bundle) Sum() string {
	hasher := sha256.New()

	updateHashWithRegoFiles(hasher, b.RegoFiles)
	updateHashWithOtherFiles(hasher, b.OtherFiles)

	return hex.EncodeToString(hasher.Sum(nil))
}

func updateHashWithRegoFiles(hasher hash.Hash, files map[string]*regofile.File) {
	for _, k := range sortedMap(files) {
		hasher.Write([]byte(files[k].Sum()))
	}
}

func updateHashWithOtherFiles(hasher hash.Hash, files map[string][]byte) {
	for _, k := range sortedMap(files) {
		hasher.Write(files[k])
	}
}

func sortedMap[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

func (b *Bundle) SetRequire(requirement *Bundle, direction lockfile.DirectionType) error {
	b.requireCache.Put(requirement.Repository(), requirement.Name())
	b.LockFile.Require.List = append(b.LockFile.Require.List, &lockfile.RequirementDecl{
		Repository: requirement.Repository(),
		Direction:  direction,
		Name:       requirement.Name(),
		Version:    requirement.Version.String(),
		H1:         requirement.BundleFile.Sum(),
		H2:         requirement.Sum(),
	})

	b.BundleFile.Require.List = append(b.BundleFile.Require.List, &bundlefile.RequirementDecl{
		Repository: requirement.Repository(),
		Name:       requirement.Name(),
		Version:    requirement.Version.String(),
	})

	b.LockFile.Sum = b.BundleFile.Sum()
	return nil
}

// TODO: get rid
func (b *Bundle) Configure() *Bundle {
	b.requireCache = hashbidimap.New[string, string]()

	if b.BundleFile.Require == nil {
		b.BundleFile.Require = &bundlefile.RequireDecl{
			List: make([]*bundlefile.RequirementDecl, 0),
		}
	}

	if b.LockFile.Require == nil {
		b.LockFile.Require = &lockfile.RequireDecl{
			List: make([]*lockfile.RequirementDecl, 0),
		}
	}

	return b
}
