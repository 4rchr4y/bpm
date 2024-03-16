package storageiface

import "github.com/4rchr4y/bpm/bundle"

type Storage interface {
	Some(repo string, version string) bool
	StoreSome(b *bundle.Bundle) error
	Store(b *bundle.Bundle) error
	Load(source string, version *bundle.VersionSpec) (*bundle.Bundle, error)
	LoadFromAbs(path string, v *bundle.VersionSpec) (*bundle.Bundle, error)
}
