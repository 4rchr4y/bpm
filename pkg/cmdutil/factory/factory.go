package factory

import (
	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/internal/gitfacade"
	"github.com/4rchr4y/bpm/pkg/bundleutil/encode"
	"github.com/4rchr4y/bpm/pkg/bundleutil/inspect"
	"github.com/4rchr4y/bpm/pkg/bundleutil/manifest"
	"github.com/4rchr4y/bpm/pkg/download"
	"github.com/4rchr4y/bpm/pkg/fetch"
	"github.com/4rchr4y/bpm/pkg/storage"
	"github.com/4rchr4y/godevkit/v3/syswrap/ioiface"
	"github.com/4rchr4y/godevkit/v3/syswrap/osiface"
)

type Factory struct {
	Name    string // executable name
	Version string // app version
	Dir     string

	IOStream  core.IO
	Storage   *storage.Storage
	Inspector *inspect.Inspector
	Encoder   *encode.Encoder
	// OsLoader   *osload.OsLoader       // bundle file loader from file system
	GitFacade *gitfacade.GitFacade // facade for interaction with both the CLI and the API
	// GitLoader  *gitload.GitLoader     // bundle file loader from the git repo
	Fetcher    *fetch.Fetcher
	Downloader *download.Downloader // downloader of a bundle and its dependencies
	Manifester *manifest.Manifester // bundle manifest file control operator
	OS         osiface.OSWrapper    // set of functions for working with the OS
	IO         ioiface.IOWrapper    // set of functions for working with input/output
}
