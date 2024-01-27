package init

import (
	"fmt"
	"path"

	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/bpm/pkg/bundle/bundlefile"
	"github.com/4rchr4y/bpm/pkg/bundle/lockfile"
	"github.com/4rchr4y/bpm/pkg/encode"
)

func bundleFileContent(encoder *encode.BundleEncoder, repo string, author *bundle.AuthorExpr) []byte {
	repoName := path.Base(repo)
	bundlefile := &bundlefile.File{
		Package: &bundlefile.PackageDecl{
			Name:        repoName,
			Author:      []string{author.String()},
			Repository:  repo,
			Description: fmt.Sprintf("Some description about '%s' bundle.", repoName),
		},
	}

	return encoder.EncodeBundleFile(bundlefile)
}

func lockfileContent(encoder *encode.BundleEncoder, sum string, edition string) []byte {
	lockfile := &lockfile.File{
		Sum:     sum,
		Edition: edition,
	}

	return encoder.EncodeLockFile(lockfile)
}

func gitignoreFileContent() []byte {
	return []byte(`bundle.lock`)
}

func bpmignoreFileContent() []byte {
	return []byte(`.git`)
}
