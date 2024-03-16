package bpmiface

import "github.com/4rchr4y/bpm/bundle"

type Inspector interface {
	Inspect(b *bundle.Bundle) error
	Verify(b *bundle.Bundle) (err error)
	Validate(b *bundle.Bundle) (err error)
}
