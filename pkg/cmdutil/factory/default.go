package factory

import (
	"github.com/4rchr4y/bpm/pkg/bundleutil"
	"github.com/4rchr4y/bpm/pkg/iostream"
	"github.com/4rchr4y/godevkit/syswrap"

	"github.com/4rchr4y/bpm/internal/gitfacade"
)

func New(version string) *Factory {
	io := iostream.NewIOStream()

	osWrap := new(syswrap.OsWrapper)
	ioWrap := new(syswrap.IoWrapper)
	encoder := bundleutil.NewEncoder()

	manifester := bundleutil.NewManifester(io, osWrap, encoder)
	fileifier := bundleutil.NewFileifier(io, encoder, manifester)
	verifier := bundleutil.NewVerifier(io)

	gitFacade := gitfacade.NewGitFacade()
	loader := bundleutil.NewLoader(io, osWrap, ioWrap, fileifier, gitFacade)

	f := &Factory{
		Name:      "bpm",
		Version:   version,
		IOStream:  io,
		Encoder:   encoder,
		Fileifier: fileifier,
		Loader:    loader,
		GitFacade: gitFacade,
		// GitLoader:  loader,
		Saver:      bundleutil.NewSaver(osWrap, encoder),
		Downloader: bundleutil.NewDownloader(loader, verifier),
		Manifester: manifester,
		Verifier:   verifier,
		IO:         ioWrap,
		OS:         osWrap,
	}

	return f
}
