package bundleutil

import "github.com/4rchr4y/bpm/pkg/bundle"

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) Validate(b *bundle.Bundle) error {
	return nil
}
