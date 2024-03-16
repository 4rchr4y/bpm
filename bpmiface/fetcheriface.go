package bpmiface

import (
	"context"

	"github.com/4rchr4y/bpm/bundle"
	"github.com/4rchr4y/bpm/fetch"
)

type Fetcher interface {
	Fetch(ctx context.Context, source string, version *bundle.VersionSpec) (*fetch.FetchOutput, error)
	PlainFetch(ctx context.Context, source string, version *bundle.VersionSpec) (*bundle.Bundle, error)
	FetchLocal(ctx context.Context, source string, version *bundle.VersionSpec) (*bundle.Bundle, error)
	FetchRemote(ctx context.Context, source string, version *bundle.VersionSpec) (*bundle.Bundle, error)
}
