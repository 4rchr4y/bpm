package factory

import (
	"github.com/4rchr4y/bpm/pkg/bundleutil"
	"github.com/4rchr4y/bpm/pkg/load/gitload"
	"github.com/4rchr4y/bpm/pkg/load/osload"
	"github.com/4rchr4y/godevkit/syswrap"

	"github.com/4rchr4y/bpm/internal/gitfacade"
)

func New(version string) *Factory {
	osWrap := new(syswrap.OsWrapper)
	ioWrap := new(syswrap.IoWrapper)
	encoder := bundleutil.NewEncoder()

	manifester := bundleutil.NewManifester(osWrap, encoder)
	fileifier := bundleutil.NewFileifier(encoder, manifester)
	verifier := bundleutil.NewVerifier()

	gitFacade := gitfacade.NewGitFacade()
	gitLoader := gitload.NewGitLoader(gitFacade, fileifier)
	osLoader := osload.NewOsLoader(osWrap, ioWrap, fileifier)

	f := &Factory{
		Name:       "bpm",
		Version:    version,
		Encoder:    encoder,
		Fileifier:  fileifier,
		OsLoader:   osLoader,
		GitFacade:  gitFacade,
		GitLoader:  gitLoader,
		Saver:      bundleutil.NewSaver(osWrap, encoder),
		Downloader: bundleutil.NewDownloader(gitLoader, verifier),
		Manifester: manifester,
		Verifier:   verifier,
		IO:         ioWrap,
		OS:         osWrap,
	}

	return f
}
