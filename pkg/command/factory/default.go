package factory

import (
	"github.com/4rchr4y/bpm/pkg/bundleutil"
	"github.com/4rchr4y/bpm/pkg/fileify"
	"github.com/4rchr4y/bpm/pkg/load/gitload"
	"github.com/4rchr4y/bpm/pkg/load/osload"
	"github.com/4rchr4y/godevkit/syswrap"

	gitcli "github.com/4rchr4y/bpm/internal/git"
)

func New(version string) *Factory {
	osWrap := new(syswrap.OsWrapper)
	ioWrap := new(syswrap.IoWrapper)
	encoder := bundleutil.NewEncoder()
	fileifier := fileify.NewFileifier(encoder)

	gitLoader := gitload.NewGitLoader(gitcli.NewClient(), fileifier)
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
