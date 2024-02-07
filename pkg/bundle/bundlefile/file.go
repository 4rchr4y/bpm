package bundlefile

import (
	"crypto/sha256"
	"fmt"

	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/pkg/util"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/samber/lo"
)

type AuthorExpr struct {
	Username string // value of git 'config --get user.username'
	Email    string // value of git 'config --get user.email'
}

func (author *AuthorExpr) String() string {
	return fmt.Sprintf("%s %s", author.Username, author.Email)
}

type PackageDecl struct {
	Name        string   `hcl:"name"`
	Author      []string `hcl:"author,optional"`
	Repository  string   `hcl:"repository"`
	Description string   `hcl:"description,optional"`
}

type RequirementDecl struct {
	Repository string `hcl:"repository,label"`
	Name       string `hcl:"name"`
	Version    string `hcl:"version"`
}

type RequireDecl struct {
	List []*RequirementDecl `hcl:"bundle,block"`
}

type File struct {
	Package *PackageDecl `hcl:"package,block"`
	Require *RequireDecl `hcl:"require,block"`
}

func (*File) Filename() string { return constant.BundleFileName }

func (bf *File) Sum() string {
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(bf, f.Body())
	return util.ChecksumSHA256(sha256.New(), f.Bytes())
}

func (bf *File) IsRequirementListed(repo string, version string) bool {
	if bf.Require == nil {
		return false
	}

	return lo.SomeBy(bf.Require.List, func(item *RequirementDecl) bool {
		if item.Repository == repo {
			return item.Version == version
		}

		return false
	})
}
