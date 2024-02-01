package lockfile

import (
	"github.com/4rchr4y/bpm/constant"
)

type DirectionType string

const (
	Direct   = "direct"
	Indirect = "indirect"
)

func (dt DirectionType) String() string {
	return string(dt)
}

type (
	ModDecl struct {
		Package string   `hcl:"package,label"` // rego file package name, 		e.g. 'data.example'
		Source  string   `hcl:"source"`        // file source path, 			e.g. 'example/file.rego'
		Sum     string   `hcl:"sum"`           // calculated file checksum		e.g. 'd973b71fd6dd925...'
		Require []string `hcl:"require"`       // direct module dependencies 	e.g. '[...]'
	}

	ModulesDecl struct {
		List []*ModDecl `hcl:"mod,block"`
	}
)

type (
	RequirementDecl struct {
		Repository string `hcl:"repository,label"` // bundle repository url						e.g. 'github.com/4rchr4y/example'
		Direction  string `hcl:"direction,label"`  // direction type, e.g. direct or indirect	e.g. 'direct'
		Name       string `hcl:"name"`             // name form bundle file						e.g. 'example'
		Version    string `hcl:"version"`          // bundle version							e.g. 'v0.0.0+20240128102927-ab4647768668'
		H1         string `hcl:"h1"`               // bundle file checksum						e.g. 'd973b71fd6dd925...'
		H2         string `hcl:"h2"`               // bundle files + other files checksum		e.g. 'd973b71fd6dd925...'
	}

	RequireDecl struct {
		List []*RequirementDecl `hcl:"bundle,block"`
	}
)

type File struct {
	Sum     string       `hcl:"sum"`           // bundle file checksum				e.g. 'd973b71fd6dd925...'
	Edition string       `hcl:"edition"`       // lock file edition 				e.g. '2024'
	Modules *ModulesDecl `hcl:"modules,block"` // list of nested modules
	Require *RequireDecl `hcl:"require,block"` // list of declared dependencies
}

func (*File) Filename() string { return constant.LockFileName }
