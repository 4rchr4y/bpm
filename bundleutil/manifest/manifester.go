package manifest

import (
	"fmt"
	"path/filepath"

	"github.com/4rchr4y/bpm/bundle"
	"github.com/4rchr4y/bpm/bundle/bundlefile"
	"github.com/4rchr4y/bpm/bundle/lockfile"
	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/core"
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
	EncodeBundleFile(bundlefile *bundlefile.Schema) []byte
	EncodeLockFile(lockfile *lockfile.Schema) []byte
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
			lockfile.FilterBySource(requirement.Repository()),
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

var compOp = map[bool]string{true: "=>", false: "<="}

func (m *Manifester) createRdirect(parent *bundle.Bundle, comingRequirement *bundle.Bundle) error {
	ebr, ebrIdx, ebrOk := parent.BundleFile.FindIndexOfRequirement(
		bundlefile.FilterBySource(comingRequirement.Repository()),
		bundlefile.FilterByVersion(comingRequirement.Version.String()),
	)
	_, _, elrOk := parent.LockFile.FindIndexOfRequirement(
		lockfile.FilterBySource(comingRequirement.Repository()),
		lockfile.FilterByVersion(comingRequirement.Version.String()),
	)

	if ebrOk && elrOk {
		m.io.PrintfOk("bundle %s is already updated", bundle.FormatVersionFromBundle(comingRequirement))
		return nil
	}

	if !elrOk {
		parent.LockFile.Require.List = append(
			parent.LockFile.Require.List,
			NewLockfileRequirementDecl(comingRequirement, lockfile.Direct),
		)
	}

	if !ebrOk {
		parent.BundleFile.Require.List = append(
			parent.BundleFile.Require.List,
			NewBundlefileRequirementDecl(comingRequirement),
		)
	}

	if ebr != nil {
		actualVersion, err := bundle.ParseVersionExpr(ebr.Version)
		if err != nil {
			return nil
		}

		isGreater := actualVersion.GreaterThan(comingRequirement.Version)
		if !isGreater {
			m.io.PrintfWarn("installing an older bundle %s version", ebr.Repository)
		}

		m.io.PrintfInfo("upgrading %s %s %s",
			bundle.FormatVersion(ebr.Repository, ebr.Version),
			compOp[isGreater],
			bundle.FormatVersionFromBundle(comingRequirement),
		)

		parent.BundleFile.Require.List[ebrIdx] = NewBundlefileRequirementDecl(comingRequirement)
	}

	m.io.PrintfOk("bundle %s has been successfully added", bundle.FormatVersionFromBundle(comingRequirement))
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
