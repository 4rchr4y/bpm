package factory

import (
	"github.com/4rchr4y/bpm/internal/service/github"
	"github.com/4rchr4y/bpm/pkg/bundleutil/encode"
	"github.com/4rchr4y/bpm/pkg/bundleutil/inspect"
	"github.com/4rchr4y/bpm/pkg/bundleutil/manifest"
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
	encoder := &encode.Encoder{
		IO: io,
	}

	manifester := manifest.NewManifester(io, osWrap, encoder)
	inspector := &inspect.Inspector{
		IO: io,
	}

	githubClient := &github.GitClient{}
	githubCLI := &github.GitCLI{}

	storage := &storage.Storage{
		Dir:     dir,
		IO:      io,
		OSWrap:  osWrap,
		IOWrap:  ioWrap,
		Encoder: encoder,
	}

	fetcher := &fetch.Fetcher{
		IO:           io,
		Storage:      storage,
		Inspector:    inspector,
		GitHubClient: githubClient,
		Encoder:      encoder,
	}

	f := &Factory{
		Name:         "bpm",
		Version:      version,
		IOStream:     io,
		Encoder:      encoder,
		Inspector:    inspector,
		Fetcher:      fetcher,
		Storage:      storage,
		GitHubClient: githubClient,
		GitCLI:       githubCLI,
		Manifester:   manifester,
		IO:           ioWrap,
		OS:           osWrap,
	}

	return f
}
