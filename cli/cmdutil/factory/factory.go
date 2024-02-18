package factory

import (
	"github.com/4rchr4y/bpm/bundleutil/encode"
	"github.com/4rchr4y/bpm/bundleutil/inspect"
	"github.com/4rchr4y/bpm/bundleutil/manifest"
	"github.com/4rchr4y/bpm/fetch"
	"github.com/4rchr4y/bpm/internal/service/github"
	"github.com/4rchr4y/bpm/iostream/iostreamiface"
	"github.com/4rchr4y/bpm/storage"
	"github.com/4rchr4y/godevkit/v3/syswrap/ioiface"
	"github.com/4rchr4y/godevkit/v3/syswrap/osiface"
)

type Factory struct {
	Name    string // executable name
	Version string // app version
	Dir     string

	IOStream   iostreamiface.IO
	Storage    *storage.Storage
	Inspector  *inspect.Inspector
	Encoder    *encode.Encoder
	GitCLI     *github.GitCLI
	Fetcher    *fetch.Fetcher
	Manifester *manifest.Manifester // bundle manifest file control operator
	OS         osiface.OSWrapper    // set of functions for working with the OS
	IO         ioiface.IOWrapper    // set of functions for working with input/output
}
