package bundle

import "github.com/4rchr4y/bpm/constant"

type Requirement struct {
	Package    string       `toml:"version" validate:"required"`
	Version    *VersionExpr `toml:"version" validate:"required"`
	Repository string       `toml:"repository" validate:"required"`
}

type ModuleDef struct {
	Name     string         `toml:"name"`
	Source   string         `toml:"source"`
	Checksum string         `toml:"checksum"`
	Require  []*Requirement `toml:"require"`
}

type BundleLockFile struct {
	Version int          `toml:"version"` // file revision version
	Modules []*ModuleDef `toml:"modules"`
}

func (*BundleLockFile) bpmFile()         {}
func (*BundleLockFile) FileName() string { return constant.LockFileName }
