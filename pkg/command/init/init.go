package init

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os/exec"
	"strings"

	"github.com/4rchr4y/bpm/cli/require"
	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/bpm/pkg/command/factory"
	"github.com/4rchr4y/bpm/pkg/encode"
	"github.com/spf13/cobra"
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

func NewCmdInit(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init NAME",
		Aliases: []string{"new", "create"},
		Args:    require.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				// Allow file completion when completing the argument for the name
				// which could be a path
				return nil, cobra.ShellCompDirectiveDefault
			}

			// No more completions, so disable file completion
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Short: "Init a new bundle.",
		Long:  cmdInitDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			return initRun(&initOptions{
				Repository: args[0],
				Author: func() *bundle.AuthorExpr {
					username, _ := getGitUserInfo("username")
					email, _ := getGitUserInfo("email")

					if username == "" || email == "" {
						return nil
					}

					return &bundle.AuthorExpr{
						Username: username,
						Email:    email,
					}
				}(),
				Encoder:   f.Encoder,
				WriteFile: f.OS.WriteFile,
			})
		},
	}

	return cmd
}

type initOptions struct {
	Repository string                                                 // repo to which the bundle will belong
	Author     *bundle.AuthorExpr                                     // git information about the author
	Encoder    *encode.BundleEncoder                                  // decoder of bundle component files
	WriteFile  func(name string, data []byte, perm fs.FileMode) error // func of saving a file to disk
}

func initRun(opts *initOptions) error {
	bundlefileContent := bundleFileContent(opts.Encoder, opts.Repository, opts.Author)
	bundlefileHash := md5.Sum(bytes.TrimSpace(bundlefileContent))
	lockfileContent := lockfileContent(opts.Encoder, hex.EncodeToString(bundlefileHash[:]), "2024")

	files := map[string][]byte{
		".gitignore":            gitignoreFileContent(),
		constant.BundleFileName: bundlefileContent,
		constant.LockFileName:   lockfileContent,
		constant.IgnoreFileName: bpmignoreFileContent(),
	}

	for fileName, content := range files {
		if err := opts.WriteFile(fileName, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file '%s': %v", fileName, err)
		}
	}

	return nil
}

func getGitUserInfo(target string) (string, error) {
	cmd := exec.Command("git", "config", "--get", fmt.Sprintf("user.%s", target))
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
