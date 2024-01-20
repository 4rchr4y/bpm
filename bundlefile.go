package bpm

import "fmt"

type BPMFile interface {
	Name() string

	bpmFile()
}

// ----------------- Bundle File ----------------- //

type PackageDef struct {
	Name        string   `toml:"name" validate:"required"`
	Version     string   `toml:"version" validate:"required"`
	Author      []string `toml:"author"`
	Description string   `toml:"description"`
}

type DependencyDef struct {
	Version string   `toml:"version"`
	Source  string   `toml:"source"`
	Include []string `toml:"include"`
}

type BundleFile struct {
	Package      *PackageDef               `toml:"package" validate:"required"`
	Dependencies map[string]*DependencyDef `toml:"dependencies"`
}

func (*BundleFile) bpmFile()     {}
func (*BundleFile) Name() string { return BPMBundleFile }

type validateClient interface {
	ValidateStruct(s interface{}) error
}

func (bf *BundleFile) Validate(validator validateClient) error {
	if err := validator.ValidateStruct(bf); err != nil {
		return fmt.Errorf("failed to validate %s file: %v", bf.Name(), err)
	}

	return nil
}

// ----------------- Bundle Lock File ----------------- //

type ModuleDef struct {
	Name         string   `toml:"name"`
	Source       string   `toml:"source"`
	Checksum     string   `toml:"checksum"`
	Dependencies []string `toml:"dependencies"`
}

type BundleLockFile struct {
	Version int          `toml:"version"`
	Modules []*ModuleDef `toml:"modules"`
}

func (*BundleLockFile) bpmFile()     {}
func (*BundleLockFile) Name() string { return BPMLockFile }

// ----------------- Bundle Work File ----------------- //

type WorkspaceDef struct {
	Path     string   `toml:"path"`
	Author   []string `toml:"author"`
	Packages []string `toml:"packages"`
}

type BpmWorkFile struct {
	Workspace *WorkspaceDef `toml:"workspace"`
}

func (*BpmWorkFile) bpmFile()     {}
func (*BpmWorkFile) Name() string { return BPMWorkFile }
