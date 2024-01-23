package bundle

import (
	"fmt"

	"github.com/4rchr4y/bpm/constant"
)

type BundleFilePackage struct {
	Name        string   `toml:"name" validate:"required"`
	Author      []string `toml:"author,omitempty"`
	Repository  string   `toml:"repository" validate:"required"`
	Description string   `toml:"description,omitempty"`
}

type BundleFileRequirement struct {
	Name    string `toml:"name" validate:"required"`
	Version string `toml:"version" validate:"required"`
}

type BundleFile struct {
	Package *BundleFilePackage                `toml:"package" validate:"required"`
	Require map[string]*BundleFileRequirement `toml:"require,omitempty"`
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
