package init

import (
	"fmt"
	"path"

	"github.com/4rchr4y/bpm/bundle/bundlefile"
	"github.com/4rchr4y/bpm/bundle/lockfile"
	"github.com/4rchr4y/bpm/bundleutil/encode"
)

const cmdInitDesc = `
The 'bpm init' command is designed to generate a new bundle directory 
structure, complete with standard files typical for a bundle.

For instance, executing bpm new foo would result in a directory 
setup similar to the following:

foo/
├── .bpmignore		# Ignore files for bpm.
├── .gitignore		# Ignore for git system.
└── bundle.hcl		# File with bundle information.

'bpm init' takes a path for an argument. If directories in the given path
do not exist, bpm will attempt to create them as it goes. If the given
destination exists and there are files in that directory, conflicting files
will be overwritten, but other files will be left alone.
`

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
