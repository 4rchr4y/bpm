package bundleutil

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/bpm/pkg/bundle/bundlefile"
	"github.com/4rchr4y/bpm/pkg/bundle/lockfile"
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

type UpdateInput struct {
	Target    *bundle.Bundle   // target bundle that needed to be updated
	Rdirect   []*bundle.Bundle // directly incoming required bundles
	Rindirect []*bundle.Bundle // indirectly incoming required bundles
}

func (m *Manifester) Update(input *UpdateInput) error {
	if input.Target.BundleFile.Require == nil && input.Rdirect != nil {
		input.Target.BundleFile.Require = &bundlefile.RequireDecl{
			List: make([]*bundlefile.RequirementDecl, len(input.Rdirect)),
		}
	}

	if input.Target.LockFile.Require == nil && input.Rindirect != nil {
		input.Target.LockFile.Require = &lockfile.RequireDecl{
			List: make([]*lockfile.RequirementDecl, len(input.Rindirect)),
		}
	}

	for _, requirement := range input.Rdirect {
		input.Target.BundleFile.Require.List = append(input.Target.BundleFile.Require.List, &bundlefile.RequirementDecl{
			Repository: requirement.Repository(),
			Name:       requirement.Name(),
			Version:    requirement.Version.String(),
		})

		input.Target.LockFile.Require.List = append(input.Target.LockFile.Require.List, &lockfile.RequirementDecl{
			Repository: requirement.Repository(),
			Direction:  lockfile.Direct,
			Name:       requirement.Name(),
			Version:    requirement.Version.String(),
			H1:         requirement.BundleFile.Sum(),
			H2:         requirement.Sum(),
		})
	}

	for _, requirement := range input.Rindirect {
		input.Target.LockFile.Require.List = append(input.Target.LockFile.Require.List, &lockfile.RequirementDecl{
			Repository: requirement.Repository(),
			Direction:  lockfile.Indirect,
			Name:       requirement.Name(),
			Version:    requirement.Version.String(),
			H1:         requirement.BundleFile.Sum(),
			H2:         requirement.Sum(),
		})
	}

	input.Target.LockFile.Sum = input.Target.Sum()

	m.io.PrintfDebug("bundle '%s' has been successfully updated", input.Target.Repository())
	return nil
}

func (m *Manifester) Upgrade(workDir string, b *bundle.Bundle) error {
	if err := m.upgradeBundleFile(workDir, b); err != nil {
		return err
	}

	if err := m.upgradeLockFile(workDir, b); err != nil {
		return err
	}

	m.io.PrintfDebug("bundle '%s' has been successfully upgraded", b.Repository())
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
