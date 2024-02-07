package factory

import (
	"github.com/4rchr4y/bpm/pkg/bundleutil"
	"github.com/4rchr4y/bpm/pkg/fetch"
	"github.com/4rchr4y/bpm/pkg/iostream"
	"github.com/4rchr4y/godevkit/syswrap"

	"github.com/4rchr4y/bpm/internal/gitfacade"
)

func New(version string) *Factory {
	io := iostream.NewIOStream()

	osWrap := new(syswrap.OSWrap)
	ioWrap := new(syswrap.IOWrap)
	encoder := bundleutil.NewEncoder()

	manifester := bundleutil.NewManifester(io, osWrap, encoder)
	fileifier := bundleutil.NewFileifier(io, encoder, manifester)
	verifier := bundleutil.NewVerifier(io)

	gitFacade := gitfacade.NewGitFacade()

	fetcher := &fetch.Fetcher{
		IO:        io,
		OSWrap:    osWrap,
		IOWrap:    ioWrap,
		Verifier:  verifier,
		Fileifier: fileifier,
		GitFacade: gitFacade,
	}

	f := &Factory{
		Name:      "bpm",
		Version:   version,
		IOStream:  io,
		Encoder:   encoder,
		Fileifier: fileifier,
		Fetcher:   fetcher,
		GitFacade: gitFacade,
		// GitLoader:  loader,
		Saver:      bundleutil.NewSaver(osWrap, encoder),
		Downloader: bundleutil.NewDownloader(fetcher, verifier),
		Manifester: manifester,
		Verifier:   verifier,
		IO:         ioWrap,
		OS:         osWrap,
	}

	return f
}
