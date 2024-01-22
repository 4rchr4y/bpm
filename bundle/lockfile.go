package bundle

import "github.com/4rchr4y/bpm/constant"

type ModuleDef struct {
	Name         string   `toml:"name"`
	Source       string   `toml:"source"`
	Checksum     string   `toml:"checksum"`
	Dependencies []string `toml:"dependencies"`
}

type DependencyDef struct {
	Version string   `toml:"version"`
	Source  string   `toml:"source"`
	Include []string `toml:"include"`
}

type BundleLockFile struct {
	Version int          `toml:"version"`
	Modules []*ModuleDef `toml:"modules"`
}

func (*BundleLockFile) bpmFile()         {}
func (*BundleLockFile) FileName() string { return constant.LockFileName }
