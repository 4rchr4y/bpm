package loader

import "github.com/4rchr4y/bpm/bundle"

type bundleFileifier interface {
	Fileify(files map[string][]byte) (*bundle.Bundle, error)
}
