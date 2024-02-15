package inspect

import (
	"errors"
	"fmt"
	"strings"

	"github.com/4rchr4y/bpm/constant"
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

type VerificationError struct {
	error
}

func (e VerificationError) Error() string {
	return fmt.Sprintf("verification: %v", e.error)
}

func (insp *Inspector) Verify(b *bundle.Bundle) error {
	if b.LockFile.Sum != b.Sum() {
		return VerificationError{
			errors.New("checksum does not match the expected one"),
		}
	}

	return nil
}

type ValidationError struct {
	error
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation: %v", e.error)
}

func (insp *Inspector) Validate(b *bundle.Bundle) error {
	if b.BundleFile == nil {
		return ValidationError{
			fmt.Errorf("file %s is undefined", constant.BundleFileName),
		}
	}

	if b.LockFile == nil {
		return ValidationError{
			fmt.Errorf("file %s is undefined", constant.LockFileName),
		}
	}

	for _, f := range b.RegoFiles {
		packagePath := strings.ReplaceAll(f.Package(), ".", "/")
		sourcePath := strings.TrimSuffix(f.Path, constant.RegoFileExt)
		if packagePath != sourcePath {
			return ValidationError{
				fmt.Errorf("invalid %s package definition;\nexpected: %s\nactual: %s", f.Path, strings.ReplaceAll(sourcePath, "/", "."), f.Package()),
			}
		}
	}

	return nil
}
