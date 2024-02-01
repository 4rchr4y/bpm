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
	fileifier := bundleutil.NewFileifier(encoder)

	gitLoader := gitload.NewGitLoader(gitfacade.NewGitFacade(), fileifier)
	osLoader := osload.NewOsLoader(osWrap, ioWrap, fileifier)

	f := &Factory{
		Name:       "bpm",
		Version:    version,
		Encoder:    encoder,
		Fileifier:  fileifier,
		OsLoader:   osLoader,
		GitLoader:  gitLoader,
		Saver:      bundleutil.NewSaver(osWrap, encoder),
		Downloader: bundleutil.NewDownloader(gitLoader),
		Manifester: bundleutil.NewManifester(osWrap, encoder),
		IO:         ioWrap,
		OS:         osWrap,
	}

	return f
}
