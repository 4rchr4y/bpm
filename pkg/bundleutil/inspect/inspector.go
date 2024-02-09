package inspect

import (
	"errors"
	"fmt"

	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/pkg/bundle"
)

type Inspector struct {
	IO core.IO
}

func (insp *Inspector) Inspect(b *bundle.Bundle) error {
	if err := insp.Verify(b); err != nil {
		return fmt.Errorf("failed to process %s verification: %v", b.Repository(), err)
	}

	if err := insp.Validate(b); err != nil {
		return fmt.Errorf("failed to process %s validation: %v", b.Repository(), err)
	}

	return nil
}

func (insp *Inspector) Verify(b *bundle.Bundle) error {
	insp.IO.PrintfDebug("expected:\t%s", b.LockFile.Sum)
	insp.IO.PrintfDebug("actual:\t\t%s", b.Sum())

	if b.LockFile.Sum != b.Sum() {
		return errors.New("checksum does not match the expected one")
	}

	return nil
}

func (insp *Inspector) Validate(b *bundle.Bundle) error {
	return nil
}
