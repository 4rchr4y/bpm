package bundle

import (
	"fmt"
	"strings"

	"github.com/4rchr4y/bpm/constant"
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
	Version        *VersionExpr
	BundleFile     *BundleFile
	BundleLockFile *BundleLockFile
	IgnoreFiles    map[string]struct{}
	RegoFiles      map[string]*RawRegoFile
	OtherFiles     map[string][]byte
}

func (b *Bundle) Name() string               { return b.BundleFile.Package.Name }
func (b *Bundle) Repository() string         { return b.BundleFile.Package.Repository }
func (b *Bundle) Require() map[string]string { return b.BundleFile.Require }
