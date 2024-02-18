package lockfile

import (
	"sort"

	"github.com/4rchr4y/bpm/constant"
	"github.com/samber/lo"
)

type DirectionType string

const (
	Direct   DirectionType = "direct"
	Indirect DirectionType = "indirect"
)

var Keywords = [...]string{
	Direct.String(),
	Indirect.String(),
}

func (dt DirectionType) String() string {
	return string(dt)
}

type (
	ModDecl struct {
		Package string   `hcl:"package,label"` // rego file package name, 		e.g. 'data.example'
		Source  string   `hcl:"source"`        // file source path, 			e.g. 'example/file.rego'
		Sum     string   `hcl:"sum"`           // calculated file checksum		e.g. 'd973b71fd6dd925...'
		Require []string `hcl:"require"`       // direct module dependencies 	e.g. '{...}'
	}

	ModulesBlock struct {
		List []*ModDecl `hcl:"mod,block"` // e.g. '{...}'
	}
)

func (md *ModulesBlock) Sort() *ModulesBlock {
	sort.Slice(md.List, func(i, j int) bool {
		return md.List[i].Package < md.List[j].Package
	})

	return md
}

type (
	RequirementDecl struct {
		Repository string `hcl:"repository,label"` // bundle repository url						e.g. 'github.com/4rchr4y/example'
		Direction  string `hcl:"direction,label"`  // direction type, e.g. direct or indirect	e.g. 'direct'
		Name       string `hcl:"name"`             // name form bundle file						e.g. 'example'
		Version    string `hcl:"version"`          // bundle version							e.g. 'v0.0.0+20240128102927-ab4647768668'
		H1         string `hcl:"h1"`               // bundle file checksum						e.g. 'd973b71fd6dd925...'
		H2         string `hcl:"h2"`               // bundle files + other files checksum		e.g. 'd973b71fd6dd925...'
	}

	RequireBlock struct {
		List []*RequirementDecl `hcl:"bundle,block"`
	}
)

type Schema struct {
	Sum     string        `hcl:"sum"`           // bundle file checksum				e.g. 'd973b71fd6dd925...'
	Edition string        `hcl:"edition"`       // lock file edition 				e.g. '2024'
	Modules *ModulesBlock `hcl:"modules,block"` // list of nested modules			e.g. '{...}'
	Require *RequireBlock `hcl:"require,block"` // list of declared dependencies		e.g. '{...}'
}

func PrepareSchema(existing *Schema) *Schema {
	if existing == nil {
		return &Schema{
			Edition: "2024",
			Modules: &ModulesBlock{
				List: make([]*ModDecl, 0),
			},
			Require: &RequireBlock{
				List: make([]*RequirementDecl, 0),
			},
		}
	}

	if existing.Modules == nil {
		existing.Modules = &ModulesBlock{
			List: make([]*ModDecl, 0),
		}
	}

	if existing.Require == nil {
		existing.Require = &RequireBlock{
			List: make([]*RequirementDecl, 0),
		}
	}

	return existing
}

func (*Schema) Filename() string { return constant.LockFileName }

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

func (f *Schema) FindIndexOfRequirement(filters ...FilterFn) (*RequirementDecl, int, bool) {
	if f.Require == nil {
		return nil, -1, false
	}

	return lo.FindIndexOf(f.Require.List, func(item *RequirementDecl) bool {
		for _, filterFn := range filters {
			if !filterFn(item) {
				return false
			}
		}

		return true
	})
}
