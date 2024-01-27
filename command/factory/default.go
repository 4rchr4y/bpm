package factory

import (
	"github.com/4rchr4y/bpm/fileifier"
	"github.com/4rchr4y/bpm/pkg/encode"
	"github.com/4rchr4y/bpm/pkg/install"
	"github.com/4rchr4y/bpm/pkg/load/gitload"
	"github.com/4rchr4y/bpm/pkg/load/osload"
	"github.com/4rchr4y/godevkit/syswrap"

	gitcli "github.com/4rchr4y/bpm/internal/git"
)

func New(version string) *Factory {
	osWrap := new(syswrap.OsWrapper)
	ioWrap := new(syswrap.IoWrapper)
	encoder := encode.NewBundleEncoder()
	fileifier := fileifier.NewFileifier(encoder)

	f := &Factory{
		Name:      "bpm",
		Version:   version,
		Encoder:   encoder,
		Fileifier: fileifier,
		OsLoader:  osload.NewOsLoader(osWrap, ioWrap, fileifier),
		GitLoader: gitload.NewGitLoader(gitcli.NewClient(), fileifier),
		Installer: install.NewBundleInstaller(osWrap, encoder),

		IO: ioWrap,
		OS: osWrap,
	}

	return f
}
