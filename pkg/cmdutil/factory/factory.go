package factory

import (
	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/internal/gitfacade"
	"github.com/4rchr4y/bpm/pkg/bundleutil"
	"github.com/4rchr4y/godevkit/syswrap"
)

type Factory struct {
	Name    string // executable name
	Version string // app version

	IOStream  core.IO
	Verifier  *bundleutil.Verifier
	Encoder   *bundleutil.Encoder   // decoder of bundle component files
	Fileifier *bundleutil.Fileifier // transformer of file contents into structures
	// OsLoader   *osload.OsLoader       // bundle file loader from file system
	GitFacade *gitfacade.GitFacade // facade for interaction with both the CLI and the API
	// GitLoader  *gitload.GitLoader     // bundle file loader from the git repo
	Loader     *bundleutil.Loader
	Saver      *bundleutil.Saver      // bundle saver files into the file system
	Downloader *bundleutil.Downloader // downloader of a bundle and its dependencies
	Manifester *bundleutil.Manifester // bundle manifest file control operator
	OS         *syswrap.OsWrapper     // set of functions for working with the OS
	IO         *syswrap.IoWrapper     // set of functions for working with input/output
}
