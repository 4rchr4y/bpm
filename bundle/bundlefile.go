package bundle

import (
	"fmt"

	"github.com/4rchr4y/bpm/constant"
)

type PackageDef struct {
	Name        string   `toml:"name" validate:"required"`
	Author      []string `toml:"author"`
	Description string   `toml:"description"`
}

type DependencyDef struct {
	Version string   `toml:"version"`
	Source  string   `toml:"source"`
	Include []string `toml:"include"`
}

type WorkspaceDef struct {
	Ignore []string `toml:"ignore"`
}

type BundleFile struct {
	Package      *PackageDef               `toml:"package" validate:"required"`
	Dependencies map[string]*DependencyDef `toml:"dependencies"`
}

func (*BundleFile) bpmFile()         {}
func (*BundleFile) FileName() string { return constant.BundleFileName }

type validateClient interface {
	ValidateStruct(s interface{}) error
}

func (bf *BundleFile) Validate(validator validateClient) error {
	if err := validator.ValidateStruct(bf); err != nil {
		return fmt.Errorf("failed to validate %s file: %v", bf.FileName(), err)
	}

	return nil
}
