package init

import (
	"fmt"
	"path"

	"github.com/4rchr4y/bpm/bundle/bundlefile"
	"github.com/4rchr4y/bpm/bundle/lockfile"
	"github.com/4rchr4y/bpm/bundleutil/encode"
)

func bundleFileContent(encoder *encode.Encoder, repo string, author *bundlefile.AuthorExpr) []byte {
	repoName := path.Base(repo)
	bundlefile := &bundlefile.Schema{
		Package: &bundlefile.PackageDecl{
			Name:        repoName,
			Author:      []string{author.String()},
			Repository:  repo,
			Description: fmt.Sprintf("Some description about '%s' bundle.", repoName),
		},
	}

	return encoder.EncodeBundleFile(bundlefile)
}

func lockfileContent(encoder *encode.Encoder, sum string, edition string) []byte {
	lockfile := &lockfile.Schema{
		Sum:     sum,
		Edition: edition,
	}

	return encoder.EncodeLockFile(lockfile)
}

func bpmignoreFileContent() []byte {
	return []byte(`.git`)
}
