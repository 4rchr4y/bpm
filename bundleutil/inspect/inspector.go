package inspect

import (
	"errors"
	"fmt"
	"strings"

	"github.com/4rchr4y/bpm/bundle"
	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/core"
	"github.com/hashicorp/go-multierror"
)

type Inspector struct {
	IO core.IO
}

type VerificationError struct{ error }
type ValidationError struct{ error }

func (insp *Inspector) Inspect(b *bundle.Bundle) error {
	if err := insp.Verify(b); err != nil {
		return VerificationError{
			fmt.Errorf("failed to process %s verification: %v", b.Repository(), err),
		}
	}

	if err := insp.Validate(b); err != nil {
		return ValidationError{
			fmt.Errorf("failed to process %s validation: %v", b.Repository(), err),
		}
	}

	return nil
}

func (insp *Inspector) Verify(b *bundle.Bundle) (err error) {
	if b.LockFile.Sum != b.Sum() {
		err = multierror.Append(err, errors.New("checksum does not match the expected one"))
	}

	return err
}

func (insp *Inspector) Validate(b *bundle.Bundle) (err error) {
	if b.BundleFile == nil {
		err = multierror.Append(
			fmt.Errorf("file %s is undefined", constant.BundleFileName),
		)
	}

	if b.LockFile == nil {
		err = multierror.Append(err,
			fmt.Errorf("file %s is undefined", constant.LockFileName),
		)
	}

	for _, f := range b.RegoFiles {
		packagePath := strings.ReplaceAll(f.Package(), ".", "/")
		pathWithNoExt := strings.TrimSuffix(f.Path, constant.RegoFileExt)
		expectedPath := fmt.Sprintf("%s/%s", b.Name(), pathWithNoExt)

		if packagePath != expectedPath {
			err = multierror.Append(err, fmt.Errorf(
				"invalid %s package definition\n\t> expected: %s,\n\t> actual: %s",
				f.Path,
				strings.ReplaceAll(expectedPath, "/", "."),
				f.Package(),
			))
		}
	}

	return err
}
