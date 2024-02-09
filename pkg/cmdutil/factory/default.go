package factory

import (
	"github.com/4rchr4y/bpm/pkg/bundleutil"
	"github.com/4rchr4y/bpm/pkg/bundleutil/encode"
	"github.com/4rchr4y/bpm/pkg/bundleutil/inspect"
	"github.com/4rchr4y/bpm/pkg/bundleutil/manifest"
	"github.com/4rchr4y/bpm/pkg/fetch"
	"github.com/4rchr4y/bpm/pkg/iostream"
	"github.com/4rchr4y/bpm/pkg/storage"
	"github.com/4rchr4y/godevkit/config"
	"github.com/4rchr4y/godevkit/syswrap"

	"github.com/4rchr4y/bpm/internal/gitfacade"
)

func New() *Factory {
	dir := config.MustGetString("BPM_PATH")
	version := config.MustGetString("BPM_VERSION")

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

	fetcher := &fetch.Fetcher{
		IO:        io,
		OSWrap:    osWrap,
		IOWrap:    ioWrap,
		Inspector: inspector,
		GitFacade: gitFacade,
		Storage:   storage,
		Encoder:   encoder,
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
		Downloader: bundleutil.NewDownloader(fetcher, inspector),
		Manifester: manifester,
		IO:         ioWrap,
		OS:         osWrap,
	}

	return f
}
