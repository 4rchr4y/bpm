package lockfile

import (
	"fmt"
	"sort"

	"github.com/4rchr4y/bpm/constant"
	"github.com/samber/lo"
)

type (
	DirectionType  string
	VisibilityType string
)

const (
	Direct   DirectionType = "direct"
	Indirect DirectionType = "indirect"

	Public  VisibilityType = "public"
	Private VisibilityType = "private"
)

func (t DirectionType) String() string  { return string(t) }
func (t VisibilityType) String() string { return string(t) }

var Keywords = [...]string{
	Direct.String(),
	Indirect.String(),

	Private.String(),
	Public.String(),
}

type ModRequireSpec struct {
	Line   int
	Source string
	Module string
}

func (mre ModRequireSpec) String() string {
	return fmt.Sprintf("%d:%s:%s", mre.Line, mre.Source, mre.Module)
}

func NewModRequireSpec(line int, source, module string) ModRequireSpec {
	return ModRequireSpec{
		Line:   line,
		Source: source,
		Module: module,
	}
}

type (
	ModuleDecl struct {
		Package    string   `hcl:"package,label"`    // rego file package name, 		e.g. 'data.example'
		Visibility string   `hcl:"visibility,label"` //
		Source     string   `hcl:"source"`           // file source path, 			e.g. 'example/file.rego'
		Sum        string   `hcl:"sum"`              // calculated file checksum		e.g. 'd973b71fd6dd925...'
		Require    []string `hcl:"require,optional"` // direct module dependencies 	e.g. '{...}'
	}

	ConsistBlock struct {
		List []*ModuleDecl `hcl:"module,block"` // e.g. '{...}'
	}
)

func (md *ConsistBlock) Sort() *ConsistBlock {
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
	Consist *ConsistBlock `hcl:"consist,block"` // list of nested modules			e.g. '{...}'
	Require *RequireBlock `hcl:"require,block"` // list of declared dependencies		e.g. '{...}'
}

func PrepareSchema(existing *Schema) *Schema {
	if existing == nil {
		return &Schema{
			Edition: "2024",
			Consist: &ConsistBlock{
				List: make([]*ModuleDecl, 0),
			},
			Require: &RequireBlock{
				List: make([]*RequirementDecl, 0),
			},
		}
	}

	if existing.Consist == nil {
		existing.Consist = &ConsistBlock{
			List: make([]*ModuleDecl, 0),
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

type RequireFilterFn func(r *RequirementDecl) bool
type ModulesFilterFn func(r *ModuleDecl) bool

func ModulesFilterByPackage(pname string) ModulesFilterFn {
	return func(m *ModuleDecl) bool {
		return m.Package == pname
	}
}

func RequireFilterByVersion(version string) RequireFilterFn {
	return func(r *RequirementDecl) bool {
		return r.Version == version
	}
}

func RequireFilterBySource(source string) RequireFilterFn {
	return func(r *RequirementDecl) bool {
		return r.Repository == source
	}
}

func (bf *Schema) SomeRequirement(filters ...RequireFilterFn) bool {
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

func (f *Schema) FindIndexOfRequirement(filters ...RequireFilterFn) (*RequirementDecl, int, bool) {
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

func (s *Schema) SomeModule(filters ...ModulesFilterFn) bool {
	if s.Require == nil {
		return false
	}

	return lo.SomeBy(s.Consist.List, func(item *ModuleDecl) bool {
		for _, filterFn := range filters {
			if !filterFn(item) {
				return false
			}
		}

		return true
	})
}
