package manifest

import (
	"fmt"
	"path/filepath"

	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/bpm/pkg/bundle/bundlefile"
	"github.com/4rchr4y/bpm/pkg/bundle/lockfile"
	"github.com/4rchr4y/bpm/pkg/bundleutil"
	"github.com/4rchr4y/godevkit/v3/syswrap/osiface"
)

func NewBundlefileRequirementDecl(b *bundle.Bundle) *bundlefile.RequirementDecl {
	return &bundlefile.RequirementDecl{
		Repository: b.Repository(),
		Name:       b.Name(),
		Version:    b.Version.String(),
	}
}

func NewLockfileRequirementDecl(b *bundle.Bundle, direction lockfile.DirectionType) *lockfile.RequirementDecl {
	return &lockfile.RequirementDecl{
		Repository: b.Repository(),
		Direction:  direction.String(),
		Name:       b.Name(),
		Version:    b.Version.String(),
		H1:         b.BundleFile.Sum(),
		H2:         b.Sum(),
	}
}

type manifesterEncoder interface {
	EncodeBundleFile(bundlefile *bundlefile.File) []byte
	EncodeLockFile(lockfile *lockfile.File) []byte
}

type Manifester struct {
	io      core.IO
	osWrap  osiface.OSWrapper
	encoder manifesterEncoder
}

func NewManifester(io core.IO, osWrap osiface.OSWrapper, encoder manifesterEncoder) *Manifester {
	return &Manifester{
		io:      io,
		osWrap:  osWrap,
		encoder: encoder,
	}
}

type CreateRequirementInput struct {
	Parent    *bundle.Bundle   // target bundle that needed to be updated
	Rdirect   []*bundle.Bundle // directly incoming required bundles
	Rindirect []*bundle.Bundle // indirectly incoming required bundles
}

func (m *Manifester) CreateRequirement(input *CreateRequirementInput) error {
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

	for _, comingRequirement := range input.Rdirect {
		if err := m.createRdirect(input.Parent, comingRequirement); err != nil {
			return err
		}

	}

	for _, requirement := range input.Rindirect {
		if exists := input.Parent.LockFile.SomeRequirement(
			requirement.Repository(),
			lockfile.FilterByVersion(requirement.Version.String()),
		); !exists {
			input.Parent.LockFile.Require.List = append(
				input.Parent.LockFile.Require.List,
				NewLockfileRequirementDecl(requirement, lockfile.Indirect),
			)
		}
	}

	input.Parent.LockFile.Sum = input.Parent.Sum()

	return nil
}

func (m *Manifester) createRdirect(parent *bundle.Bundle, comingRequirement *bundle.Bundle) error {
	existingRequirement, idx, ok := parent.BundleFile.FindIndexOfRequirement(comingRequirement.Repository())
	if ok && existingRequirement.Version == comingRequirement.Version.String() {
		m.io.PrintfOk("bundle %s is already updated", bundleutil.FormatVersionFromBundle(comingRequirement))
		return nil
	}

	if existingRequirement != nil {
		input := &replaceInput{
			Parent:      parent,
			Actual:      existingRequirement,
			ActualIndex: idx,
			Coming:      comingRequirement,
		}

		if err := m.replaceBundleFileRequirement(input); err != nil {
			return err
		}

		return nil
	}

	parent.BundleFile.Require.List = append(
		parent.BundleFile.Require.List,
		NewBundlefileRequirementDecl(comingRequirement),
	)

	if exists := parent.LockFile.SomeRequirement(
		comingRequirement.Repository(),
		lockfile.FilterByVersion(comingRequirement.Version.String()),
	); !exists {
		parent.LockFile.Require.List = append(
			parent.LockFile.Require.List,
			NewLockfileRequirementDecl(comingRequirement, lockfile.Direct),
		)
	}

	m.io.PrintfOk("bundle %s has been successfully added", bundleutil.FormatVersionFromBundle(comingRequirement))
	return nil
}

type replaceInput struct {
	Parent      *bundle.Bundle
	Actual      *bundlefile.RequirementDecl
	ActualIndex int // index in the require list of the bundle file
	Coming      *bundle.Bundle
}

var compOp = map[bool]string{true: "=>", false: "<="}

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
		compOp[isGreater],
		bundleutil.FormatVersionFromBundle(input.Coming),
	)

	input.Parent.BundleFile.Require.List[input.ActualIndex] = NewBundlefileRequirementDecl(input.Coming)

	if _, idx, ok := input.Parent.LockFile.FindIndexOfRequirement(
		input.Actual.Repository,
		lockfile.FilterByVersion(input.Actual.Version),
	); ok {
		input.Parent.LockFile.Require.List[idx] = NewLockfileRequirementDecl(input.Coming, lockfile.Direct)
	}

	return nil
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
