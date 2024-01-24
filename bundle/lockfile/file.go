package lockfile

import "github.com/4rchr4y/bpm/constant"

type Requirement struct {
	Package string `toml:"version" validate:"required"`
	// Version    *VersionExpr `toml:"version" validate:"required"`
	Repository string `toml:"repository" validate:"required"`
}

type ModuleDef struct {
	Name     string         `toml:"name"`
	Source   string         `toml:"source"`
	Checksum string         `toml:"checksum"`
	Require  []*Requirement `toml:"require"`
}

type File struct {
	Version int          `toml:"version"` // file revision version
	Modules []*ModuleDef `toml:"modules"`
}

func (*File) FileName() string { return constant.LockFileName }
