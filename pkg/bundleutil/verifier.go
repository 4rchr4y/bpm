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

	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/pkg/bundle"
)

type Verifier struct {
	io core.IO
}

func NewVerifier(io core.IO) *Verifier {
	return &Verifier{
		io: io,
	}
}

func (v *Verifier) Verify(b *bundle.Bundle) error {
	v.io.PrintfDebug("expected:\t%s", b.LockFile.Sum)
	v.io.PrintfDebug("actual:\t\t%s", b.Sum())

	if b.LockFile.Sum != b.Sum() {
		return fmt.Errorf("bundle '%s' checksum does not match the expected one", b.Repository())
	}

	return nil
}
