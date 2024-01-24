package bundle

import "github.com/4rchr4y/bpm/constant"

type BundleFilePackage struct {
	Name        string   `hcl:"name"`
	Author      []string `hcl:"author,optional"`
	Repository  string   `hcl:"repository"`
	Description string   `hcl:"description,optional"`
}

type BundleFileRequirement struct {
	Name    string `hcl:"name"`
	Version string `hcl:"version"`
}

type Require struct {
	Bundle map[string]*BundleRequirement `hcl:"bundle,block"`
}

type BundleRequirement struct {
	Repository string `hcl:"repository,label"`
	Name       string `hcl:"name"`
	Version    string `hcl:"version"`
}

type RequireDecl struct {
	List []*BundleRequirement `hcl:"bundle,block"`
}

type BundleFile struct {
	Package *BundleFilePackage `hcl:"package,block"`
	Require *RequireDecl       `hcl:"require,block"`
}

func (*BundleFile) bpmFile()         {}
func (*BundleFile) FileName() string { return constant.BundleFileName }
