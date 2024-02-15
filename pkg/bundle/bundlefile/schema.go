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

type RequireBlock struct {
	List []*RequirementDecl `hcl:"bundle,block"`
}

type Schema struct {
	Package *PackageDecl  `hcl:"package,block"`
	Require *RequireBlock `hcl:"require,block"`
}

func (*Schema) Filename() string { return constant.BundleFileName }

func (bf *Schema) Sum() string {
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(bf, f.Body())
	return util.ChecksumSHA256(sha256.New(), f.Bytes())
}

type FilterFn func(r *RequirementDecl) bool

func FilterByVersion(version string) FilterFn {
	return func(r *RequirementDecl) bool {
		return r.Version == version
	}
}

func FilterBySource(source string) FilterFn {
	return func(r *RequirementDecl) bool {
		return r.Repository == source
	}
}

func FilterByName(name string) FilterFn {
	return func(r *RequirementDecl) bool {
		return r.Name == name
	}
}

func (bf *Schema) SomeRequirement(filters ...FilterFn) bool {
	if bf.Require == nil {
		return false
	}

	return lo.SomeBy(bf.Require.List, func(item *RequirementDecl) bool {
		for _, filterFn := range filters {
			if !filterFn(item) {
				return false
			}
		}

		return true
	})
}

func (bf *Schema) FindIndexOfRequirement(filters ...FilterFn) (*RequirementDecl, int, bool) {
	if bf.Require == nil {
		return nil, -1, false
	}

	return lo.FindIndexOf(bf.Require.List, func(item *RequirementDecl) bool {
		for _, filterFn := range filters {
			if !filterFn(item) {
				return false
			}
		}

		return true
	})
}
