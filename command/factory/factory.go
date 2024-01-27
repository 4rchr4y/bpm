package factory

import (
	"github.com/4rchr4y/bpm/bfencoder"
	"github.com/4rchr4y/bpm/fileifier"
	"github.com/4rchr4y/bpm/pkg/install"
	"github.com/4rchr4y/bpm/pkg/load/gitload"
	"github.com/4rchr4y/bpm/pkg/load/osload"
	"github.com/4rchr4y/godevkit/syswrap"
)

type Factory struct {
	Name    string // executable name
	Version string // app version

	Encoder   *bfencoder.Encoder       // decoder of bundle component files
	Fileifier *fileifier.Fileifier     // transformer of file contents into structures
	OsLoader  *osload.OsLoader         // bundle file loader from file system
	GitLoader *gitload.GitLoader       // bundle file loader from the git repo
	Installer *install.BundleInstaller // bundle installer into the file system
	OS        *syswrap.OsWrapper       // set of functions for working with the OS
	IO        *syswrap.IoWrapper       // set of functions for working with input/output
}
