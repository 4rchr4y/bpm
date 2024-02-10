package factory

import (
	"github.com/4rchr4y/bpm/internal/gitfacade"
	"github.com/4rchr4y/bpm/pkg/bundleutil/encode"
	"github.com/4rchr4y/bpm/pkg/bundleutil/inspect"
	"github.com/4rchr4y/bpm/pkg/bundleutil/manifest"
	"github.com/4rchr4y/bpm/pkg/download"
	"github.com/4rchr4y/bpm/pkg/fetch"
	"github.com/4rchr4y/bpm/pkg/iostream"
	"github.com/4rchr4y/bpm/pkg/storage"
	"github.com/4rchr4y/godevkit/v3/env"
	"github.com/4rchr4y/godevkit/v3/syswrap"
)

func New() *Factory {
	dir := env.MustGetString("BPM_PATH")
	version := env.MustGetString("BPM_VERSION")

	io := iostream.NewIOStream()

	osWrap := new(syswrap.OSWrap)
	ioWrap := new(syswrap.IOWrap)
	encoder := encode.NewEncoder()

	manifester := manifest.NewManifester(io, osWrap, encoder)
	inspector := &inspect.Inspector{
		IO: io,
	}

	gitFacade := gitfacade.NewGitFacade()

	storage := &storage.Storage{
		Dir:     dir,
		IO:      io,
		OSWrap:  osWrap,
		IOWrap:  ioWrap,
		Encoder: encoder,
	}

	downloader := &download.Downloader{
		IO:        io,
		OSWrap:    osWrap,
		IOWrap:    ioWrap,
		Storage:   storage,
		Inspector: inspector,
		GitFacade: gitFacade,
		Encoder:   encoder,
	}

	fetcher := &fetch.Fetcher{
		IO:         io,
		OSWrap:     osWrap,
		IOWrap:     ioWrap,
		Downloader: downloader,
		Storage:    storage,
	}

	f := &Factory{
		Name:       "bpm",
		Version:    version,
		IOStream:   io,
		Encoder:    encoder,
		Inspector:  inspector,
		Fetcher:    fetcher,
		Storage:    storage,
		GitFacade:  gitFacade,
		Downloader: downloader,
		Manifester: manifester,
		IO:         ioWrap,
		OS:         osWrap,
	}

	return f
}
