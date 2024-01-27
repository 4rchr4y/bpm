package lockfile

import (
	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/pkg/bundle/regofile"
)

type ModDecl struct {
	Package string   `hcl:"package,label"`
	Source  string   `hcl:"source"`
	Sum     string   `hcl:"sum"`
	Require []string `hcl:"require"`
}

type ModulesDecl struct {
	List []*ModDecl `hcl:"mod,block"`
}

type File struct {
	Sum     string       `hcl:"sum"`
	Edition string       `hcl:"edition"`
	Modules *ModulesDecl `hcl:"modules,block"`

	modCache map[string]struct{}
}

func (*File) FileName() string { return constant.LockFileName }

func (f *File) SetModule(file *regofile.File) {
	if f.Modules == nil {
		f.Modules = &ModulesDecl{
			List: make([]*ModDecl, 0),
		}
	}

	f.Modules.List = append(f.Modules.List, &ModDecl{
		Package: file.Package(),
		Source:  file.Path,
		Sum:     file.Sum(),
	})
}
