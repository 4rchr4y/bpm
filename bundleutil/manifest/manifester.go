package manifest

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/4rchr4y/bpm/bundle"
	"github.com/4rchr4y/bpm/bundle/bundlefile"
	"github.com/4rchr4y/bpm/bundle/lockfile"
	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/fetch"
	"github.com/4rchr4y/godevkit/v3/syswrap/osiface"
)

var compOp = map[bool]string{true: "=>", false: "<="}

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

type manifesterStorage interface {
	Store(b *bundle.Bundle) error
	Some(repo string, version string) bool
	StoreSome(b *bundle.Bundle) error
	Load(source string, version *bundle.VersionExpr) (*bundle.Bundle, error)
}

type manifesterFetcher interface {
	Fetch(ctx context.Context, source string, version *bundle.VersionExpr) (*fetch.FetchResult, error)
}

type Manifester struct {
	IO      core.IO
	OSWrap  osiface.OSWrapper
	Storage manifesterStorage
	Encoder manifesterEncoder
	Fetcher manifesterFetcher
}

type InsertRequirementInput struct {
	Parent  *bundle.Bundle
	Source  string
	Version *bundle.VersionExpr
}

func (m *Manifester) InsertRequirement(ctx context.Context, input *InsertRequirementInput) error {
	existingRequirement, idx, ok := input.Parent.BundleFile.FindIndexOfRequirement(
		bundlefile.FilterBySource(input.Source),
	)

	if ok && existingRequirement.Version == input.Version.String() {
		m.IO.PrintfOk("bundle '%s@%s' is already installed", input.Source, input.Version.String())
		return m.SyncLockfile(ctx, input.Parent) // such requirement is already installed, then just synchronize
	}

	result, err := m.Fetcher.Fetch(ctx, input.Source, input.Version)
	if err != nil {
		return err
	}

	if ok {
		existingVersion, err := bundle.ParseVersionExpr(existingRequirement.Version)
		if err != nil {
			return err
		}

		if result.Target.Version.String() == existingVersion.String() {
			m.IO.PrintfOk("bundle '%s@%s' is already installed", result.Target.Repository(), result.Target.Version.String())
			return m.SyncLockfile(ctx, input.Parent)
		}

		isGreater := result.Target.Version.GreaterThan(existingVersion)
		if !isGreater {
			m.IO.PrintfWarn("installing an older bundle %s@%s version", result.Target.Repository(), result.Target.Version.String())
		}

		m.IO.PrintfInfo("upgrading %s %s %s",
			bundle.FormatSourceVersion(input.Source, input.Version.String()),
			compOp[isGreater],
			bundle.FormatSourceVersionFromBundle(result.Target),
		)

		input.Parent.BundleFile.Require.List[idx] = NewBundlefileRequirementDecl(result.Target)

		return m.SyncLockfile(ctx, input.Parent)
	}

	input.Parent.BundleFile.Require.List = append(
		input.Parent.BundleFile.Require.List,
		NewBundlefileRequirementDecl(result.Target),
	)

	return m.SyncLockfile(ctx, input.Parent)
}

func (f *Manifester) SyncLockfile(ctx context.Context, parent *bundle.Bundle) error {
	// creating a cache for faster matching with bundle file
	requireCache := make(map[string]struct{})

	for i, req := range parent.LockFile.Require.List {
		// while going through the lock requirements, simultaneously
		// remove requirements that no longer exist in bundle file
		if exists := parent.BundleFile.SomeRequirement(
			bundlefile.FilterBySource(req.Repository),
			bundlefile.FilterByVersion(req.Version),
		); !exists {
			parent.LockFile.Require.List[i] = nil
			continue
		}

		// creating a cache of lock requirements
		formattedVersion := bundle.FormatSourceVersion(req.Repository, req.Version)
		requireCache[formattedVersion] = struct{}{}

	}

	// go through all direct requirements to ensure that the
	// lock file is up to date
	for _, r := range parent.BundleFile.Require.List {
		v, err := bundle.ParseVersionExpr(r.Version)
		if err != nil {
			return err
		}

		result, err := f.Fetcher.Fetch(ctx, r.Repository, v)
		if err != nil {
			return err
		}

		for _, b := range result.Merge() {
			if err := f.Storage.StoreSome(b); err != nil {
				return err
			}

			if _, exists := requireCache[bundle.FormatSourceVersionFromBundle(b)]; !exists {
				parent.LockFile.Require.List = append(
					parent.LockFile.Require.List,
					NewLockfileRequirementDecl(b, defineDirection(result.Target, b)),
				)
			}
		}
	}

	parent.LockFile.Sum = parent.Sum()
	return nil
}

func defineDirection(target, actual *bundle.Bundle) lockfile.DirectionType {
	if actual.Repository() == target.Repository() {
		return lockfile.Direct
	}

	return lockfile.Indirect
}

func (m *Manifester) Upgrade(workDir string, b *bundle.Bundle) error {
	if err := m.upgradeBundleFile(workDir, b); err != nil {
		return err
	}

	if err := m.upgradeLockFile(workDir, b); err != nil {
		return err
	}

	m.IO.PrintfDebug("bundle %s has been successfully upgraded", b.Repository())
	return nil
}

func (m *Manifester) upgradeBundleFile(workDir string, b *bundle.Bundle) error {
	bundlefilePath := filepath.Join(workDir, constant.BundleFileName)
	bytes := m.Encoder.EncodeBundleFile(b.BundleFile)

	if err := m.OSWrap.WriteFile(bundlefilePath, bytes, 0644); err != nil {
		return fmt.Errorf("error occurred while '%s' file updating: %v", constant.BundleFileName, err)
	}

	return nil
}

func (m *Manifester) upgradeLockFile(workDir string, b *bundle.Bundle) error {
	bundlefilePath := filepath.Join(workDir, constant.LockFileName)
	bytes := m.Encoder.EncodeLockFile(b.LockFile)

	if err := m.OSWrap.WriteFile(bundlefilePath, bytes, 0644); err != nil {
		return fmt.Errorf("error occurred while '%s' file updating: %v", constant.LockFileName, err)
	}

	return nil
}
