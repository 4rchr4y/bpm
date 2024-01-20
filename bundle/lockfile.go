package bundle

import "github.com/4rchr4y/bpm/constant"

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
func (*BundleLockFile) Name() string { return constant.LockFileName }
