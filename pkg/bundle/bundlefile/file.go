package bundlefile

import (
	"github.com/4rchr4y/bpm/constant"
)

type PackageDecl struct {
	Name        string   `hcl:"name"`
	Author      []string `hcl:"author,optional"`
	Repository  string   `hcl:"repository"`
	Description string   `hcl:"description,optional"`
}

type BundleRequirement struct {
	Repository string `hcl:"repository,label"`
	Name       string `hcl:"name"`
	Version    string `hcl:"version"`
}

type RequireDecl struct {
	List []*BundleRequirement `hcl:"bundle,block"`
}

type File struct {
	Package *PackageDecl `hcl:"package,block"`
	Require *RequireDecl `hcl:"require,block"`
}

func (*File) FileName() string { return constant.BundleFileName }
