package lockfile

import "github.com/4rchr4y/bpm/constant"

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
}

func (*File) FileName() string { return constant.LockFileName }
