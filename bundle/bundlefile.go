package bundle

import "github.com/4rchr4y/bpm/constant"

type BundleFilePackage struct {
	Name        string   `hcl:"name"`
	Author      []string `hcl:"author"`
	Repository  string   `hcl:"repository"`
	Description string   `hcl:"description"`
}

type BundleFileRequirement struct {
	Name    string `hcl:"name"`
	Version string `hcl:"version"`
}

type BundleFile struct {
	Package *BundleFilePackage                `hcl:"package"`
	Require map[string]*BundleFileRequirement `hcl:"require"`
}

func (*BundleFile) bpmFile()         {}
func (*BundleFile) FileName() string { return constant.BundleFileName }
