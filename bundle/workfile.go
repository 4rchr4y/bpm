package bundle

import "github.com/4rchr4y/bpm/constant"

type WorkspaceDef struct {
	Path     string   `toml:"path"`
	Author   []string `toml:"author"`
	Packages []string `toml:"packages"`
}

type BpmWorkFile struct {
	Workspace *WorkspaceDef `toml:"workspace"`
}

func (*BpmWorkFile) bpmFile()         {}
func (*BpmWorkFile) FileName() string { return constant.WorkFileName }
