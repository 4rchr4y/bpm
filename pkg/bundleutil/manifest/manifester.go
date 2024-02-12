package manifest

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/bpm/pkg/bundle/bundlefile"
	"github.com/4rchr4y/bpm/pkg/bundle/lockfile"
	"github.com/4rchr4y/bpm/pkg/bundleutil"
)

type manifesterEncoder interface {
	EncodeBundleFile(bundlefile *bundlefile.File) []byte
	EncodeLockFile(lockfile *lockfile.File) []byte
}

type manifesterOsWrapper interface {
	WriteFile(name string, data []byte, perm fs.FileMode) error
}

type Manifester struct {
	io      core.IO
	osWrap  manifesterOsWrapper
	encoder manifesterEncoder
}

func NewManifester(io core.IO, osWrap manifesterOsWrapper, encoder manifesterEncoder) *Manifester {
	return &Manifester{
		io:      io,
		osWrap:  osWrap,
		encoder: encoder,
	}
}

type AppendInput struct {
	Parent    *bundle.Bundle   // target bundle that needed to be updated
	Rdirect   []*bundle.Bundle // directly incoming required bundles
	Rindirect []*bundle.Bundle // indirectly incoming required bundles
}

func (m *Manifester) AppendBundle(input *AppendInput) error {
	if input.Parent.BundleFile.Require == nil && input.Rdirect != nil {
		input.Parent.BundleFile.Require = &bundlefile.RequireBlock{
			List: make([]*bundlefile.RequirementDecl, 0, len(input.Rdirect)),
		}
	}

	if input.Parent.LockFile.Require == nil && (input.Rindirect != nil || input.Rdirect != nil) {
		input.Parent.LockFile.Require = &lockfile.RequireBlock{
			List: make([]*lockfile.RequirementDecl, 0, len(input.Rindirect)),
		}
	}

	var wasUpdated bool

	for _, comingRequirement := range input.Rdirect {
		existingRequirement, idx, ok := input.Parent.BundleFile.FindIndexOfRequirement(comingRequirement.Repository())
		if ok && existingRequirement.Version == comingRequirement.Version.String() {
			m.io.PrintfOk("bundle %s is already updated", bundleutil.FormatVersionFromBundle(comingRequirement))
			continue
		}

		if existingRequirement != nil {
			input := &replaceInput{
				Parent:      input.Parent,
				Actual:      existingRequirement,
				ActualIndex: idx,
				Coming:      comingRequirement,
			}

			if err := m.replaceBundleFileRequirement(input); err != nil {
				return err
			}

			continue
		}

		wasUpdated = true

		m.appendBundleFileRequirement(input.Parent, comingRequirement)
	}

	for _, requirement := range input.Rindirect {
		input.Parent.LockFile.Require.List = append(input.Parent.LockFile.Require.List, &lockfile.RequirementDecl{
			Repository: requirement.Repository(),
			Direction:  lockfile.Indirect,
			Name:       requirement.Name(),
			Version:    requirement.Version.String(),
			H1:         requirement.BundleFile.Sum(),
			H2:         requirement.Sum(),
		})
	}

	if wasUpdated {
		input.Parent.LockFile.Sum = input.Parent.Sum()
	}

	return nil
}

type replaceInput struct {
	Parent      *bundle.Bundle
	Actual      *bundlefile.RequirementDecl
	ActualIndex int // index in the require list of the bundle file
	Coming      *bundle.Bundle
}

func (m *Manifester) replaceBundleFileRequirement(input *replaceInput) error {
	actualVersion, err := bundle.ParseVersionExpr(input.Actual.Version)
	if err != nil {
		return nil
	}

	isGreater := actualVersion.GreaterThan(input.Coming.Version)
	if !isGreater {
		m.io.PrintfWarn("installing an older bundle %s version", input.Actual.Repository)
	}

	m.io.PrintfInfo("upgrading %s %s %s",
		bundleutil.FormatVersion(input.Actual.Repository, input.Actual.Version),
		symbol(isGreater),
		bundleutil.FormatVersionFromBundle(input.Coming),
	)

	input.Parent.BundleFile.Require.List[input.ActualIndex] = &bundlefile.RequirementDecl{
		Name:       input.Coming.Name(),
		Repository: input.Actual.Repository,
		Version:    input.Coming.Version.String(),
	}

	// Update or add a requirement in a lockfile. Its absence is unusual
	// if present in the bundlefile but not treated as an error to handle
	// non-critical oddities gracefully.

	_, idx, ok := input.Parent.LockFile.FindIndexOfRequirement(
		input.Actual.Repository,
		lockfile.FilterByVersion(input.Actual.Version),
	)
	if !ok {
		return nil
	}

	input.Parent.LockFile.Require.List[idx] = &lockfile.RequirementDecl{
		Repository: input.Coming.Repository(),
		Direction:  lockfile.Direct,
		Name:       input.Coming.Name(),
		Version:    input.Coming.Version.String(),
		H1:         input.Coming.BundleFile.Sum(),
		H2:         input.Coming.Sum(),
	}

	return nil
}

func (m *Manifester) appendBundleFileRequirement(b *bundle.Bundle, requirement *bundle.Bundle) {
	b.BundleFile.Require.List = append(b.BundleFile.Require.List, &bundlefile.RequirementDecl{
		Repository: requirement.Repository(),
		Name:       requirement.Name(),
		Version:    requirement.Version.String(),
	})

	b.LockFile.Require.List = append(b.LockFile.Require.List, &lockfile.RequirementDecl{
		Repository: requirement.Repository(),
		Direction:  lockfile.Direct,
		Name:       requirement.Name(),
		Version:    requirement.Version.String(),
		H1:         requirement.BundleFile.Sum(),
		H2:         requirement.Sum(),
	})

	m.io.PrintfOk("bundle %s has been successfully added", bundleutil.FormatVersionFromBundle(requirement))
}

func (m *Manifester) Upgrade(workDir string, b *bundle.Bundle) error {
	if err := m.upgradeBundleFile(workDir, b); err != nil {
		return err
	}

	if err := m.upgradeLockFile(workDir, b); err != nil {
		return err
	}

	m.io.PrintfDebug("bundle %s has been successfully upgraded", b.Repository())
	return nil
}

func (m *Manifester) upgradeBundleFile(workDir string, b *bundle.Bundle) error {
	bundlefilePath := filepath.Join(workDir, constant.BundleFileName)
	bytes := m.encoder.EncodeBundleFile(b.BundleFile)

	if err := m.osWrap.WriteFile(bundlefilePath, bytes, 0644); err != nil {
		return fmt.Errorf("error occurred while '%s' file updating: %v", constant.BundleFileName, err)
	}

	return nil
}

func (m *Manifester) upgradeLockFile(workDir string, b *bundle.Bundle) error {
	bundlefilePath := filepath.Join(workDir, constant.LockFileName)
	bytes := m.encoder.EncodeLockFile(b.LockFile)

	if err := m.osWrap.WriteFile(bundlefilePath, bytes, 0644); err != nil {
		return fmt.Errorf("error occurred while '%s' file updating: %v", constant.LockFileName, err)
	}

	return nil
}

func symbol(isGreater bool) string {
	if isGreater {
		return "=>"
	} else {
		return "<="
	}
}
