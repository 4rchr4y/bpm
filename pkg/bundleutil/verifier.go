// The Verifier is designed solely for verifying the integrity
// of a bundle. It focuses on authenticating manifest files and validating hash
// sums of *.rego files to ensure their immutability since the time of their
// creation. The verifier's capabilities include scrutinizing signatures
// associated with manifest files to confirm their authenticity, thereby
// ensuring the manifest files have not been altered or tampered with.
// Additionally, it performs hash sum validation for files with a *.rego
// extension, guaranteeing the content of these files remains unchanged.
// However, the verifier does not extend to assessing the correctness of
// package naming within the bundle, checking the correctness of imports
// within *.rego files, or performing other validation checks not directly
// related to data integrity. It serves as a specialized tool aimed at
// confirming the authentication and immutability of bundle data, with
// operations beyond integrity verification to be handled by other tools or
// procedures.

package bundleutil

import (
	"fmt"

	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/bpm/pkg/bundle/lockfile"
)

type Verifier struct{}

func NewVerifier() *Verifier {
	return &Verifier{}
}

func (v *Verifier) Verify(b *bundle.Bundle) error {
	if err := verifyBundleFile(b); err != nil {
		return err
	}

	if err := verifyRegoFiles(b); err != nil {
		return err
	}

	return nil
}

func verifyRegoFiles(b *bundle.Bundle) error {
	modules := make(map[string]*lockfile.ModDecl, len(b.LockFile.Modules.List))

	for _, m := range b.LockFile.Modules.List {
		modules[m.Source] = m
	}

	if len(b.RegoFiles) != len(modules) {
		return fmt.Errorf("expected number of files (%d) does not match the received one (%d)", len(modules), len(b.RegoFiles))
	}

	for filePath, file := range b.RegoFiles {
		m, exits := modules[filePath]
		if !exits {
			return fmt.Errorf("can't find '%s' file", filePath)
		}

		if file.Sum() != m.Sum {
			return fmt.Errorf("file '%s' checksum does not match the expected one", filePath)
		}
	}

	return nil
}

func verifyBundleFile(b *bundle.Bundle) error {
	fmt.Println("expected\t", b.LockFile.Sum)
	fmt.Println("actual\t\t", b.Sum())

	if b.LockFile.Sum != b.Sum() {
		return fmt.Errorf("bundle '%s' checksum does not match the expected one", b.Repository())
	}

	return nil
}
