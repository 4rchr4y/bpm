package inspect

import (
	"errors"
	"fmt"
	"strings"

	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/hashicorp/go-multierror"
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

type VerificationError struct{ error }

func (e VerificationError) Error() string { return e.Error() }

func (insp *Inspector) Verify(b *bundle.Bundle) (err error) {
	if b.LockFile.Sum != b.Sum() {
		err = multierror.Append(err, errors.New("checksum does not match the expected one"))
	}

	return err
}

type ValidationError struct{ error }

func (e ValidationError) Error() string { return e.Error() }

func (insp *Inspector) Validate(b *bundle.Bundle) (err error) {
	if b.BundleFile == nil {
		err = multierror.Append(ValidationError{
			fmt.Errorf("file %s is undefined", constant.BundleFileName),
		})
	}

	if b.LockFile == nil {
		err = multierror.Append(err, ValidationError{
			fmt.Errorf("file %s is undefined", constant.LockFileName),
		})
	}

	for _, f := range b.RegoFiles {
		packagePath := strings.ReplaceAll(f.Package(), ".", "/")
		sourcePath := strings.TrimSuffix(f.Path, constant.RegoFileExt)
		if packagePath != sourcePath {
			err = multierror.Append(err, ValidationError{
				fmt.Errorf(
					"invalid %s package definition\n\t> expected: %s,\n\t> actual: %s",
					f.Path,
					strings.ReplaceAll(sourcePath, "/", "."),
					f.Package(),
				),
			})
		}
	}

	return err
}
