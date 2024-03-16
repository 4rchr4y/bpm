package bpmiface

import (
	"context"

	"github.com/4rchr4y/bpm/bundle"
	"github.com/4rchr4y/bpm/bundleutil/manifest"
)

type Manifester interface {
	SyncLockfile(ctx context.Context, parent *bundle.Bundle) error
	InsertRequirement(ctx context.Context, input *manifest.InsertRequirementInput) error
	Upgrade(workDir string, b *bundle.Bundle) error
}
