package bundlefile

import (
	"crypto/sha256"
	"fmt"

	"github.com/4rchr4y/bpm/bundleutil"
	"github.com/4rchr4y/bpm/constant"
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

type PackageBlock struct {
	Name        string   `hcl:"name"`
	Author      []string `hcl:"author,optional"`
	Repository  string   `hcl:"repository"`
	Description string   `hcl:"description,optional"`
}

type RequirementDecl struct {
	Source  string `hcl:"source,label"`
	Name    string `hcl:"name"`
	Version string `hcl:"version"`
}

type RequireBlock struct {
	List []*RequirementDecl `hcl:"bundle,block"`
}

type WorkspaceBlock struct {
	Internal []string `hcl:"internal"`
	Builtin  []string `hcl:"builtin"`
}

type Schema struct {
	Package   *PackageBlock   `hcl:"package,block"`
	Workspace *WorkspaceBlock `hcl:"workspace,block"`
	Require   *RequireBlock   `hcl:"require,block"`
}

func PrepareSchema(existing *Schema) *Schema {
	if existing == nil {
		return nil
	}

	if existing.Require == nil {
		existing.Require = &RequireBlock{
			List: make([]*RequirementDecl, 0),
		}
	}

	if existing.Workspace == nil {
		existing.Workspace = &WorkspaceBlock{
			Builtin:  make([]string, 0),
			Internal: make([]string, 0),
		}
	} else {
		if existing.Workspace.Builtin == nil {
			existing.Workspace.Builtin = make([]string, 0)
		}

		if existing.Workspace.Internal == nil {
			existing.Workspace.Internal = make([]string, 0)
		}
	}

	return existing
}

func (*Schema) Filename() string { return constant.BundleFileName }

func (s *Schema) Sum() string {
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(s, f.Body())

	return bundleutil.ChecksumSHA256(
		sha256.New(),
		bundleutil.FormatBundleFile(f.Bytes()),
	)
}

type FilterFn func(r *RequirementDecl) bool

func FilterByVersion(version string) FilterFn {
	return func(r *RequirementDecl) bool {
		return r.Version == version
	}
}

func FilterBySource(source string) FilterFn {
	return func(r *RequirementDecl) bool {
		return r.Source == source
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
