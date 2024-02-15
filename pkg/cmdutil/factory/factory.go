package factory

import (
	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/internal/service/github"
	"github.com/4rchr4y/bpm/pkg/bundleutil/encode"
	"github.com/4rchr4y/bpm/pkg/bundleutil/inspect"
	"github.com/4rchr4y/bpm/pkg/bundleutil/manifest"
	"github.com/4rchr4y/bpm/pkg/fetch"
	"github.com/4rchr4y/bpm/pkg/storage"
	"github.com/4rchr4y/godevkit/v3/syswrap/ioiface"
	"github.com/4rchr4y/godevkit/v3/syswrap/osiface"
)

type Factory struct {
	Name    string // executable name
	Version string // app version
	Dir     string

	IOStream   core.IO
	Storage    *storage.Storage
	Inspector  *inspect.Inspector
	Encoder    *encode.Encoder
	GitCLI     *github.GitCLI
	Fetcher    *fetch.Fetcher
	Manifester *manifest.Manifester // bundle manifest file control operator
	OS         osiface.OSWrapper    // set of functions for working with the OS
	IO         ioiface.IOWrapper    // set of functions for working with input/output
}
