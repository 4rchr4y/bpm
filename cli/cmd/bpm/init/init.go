package init

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/fs"

	"github.com/4rchr4y/bpm/bundle/bundlefile"
	"github.com/4rchr4y/bpm/bundleutil/encode"
	"github.com/4rchr4y/bpm/cli/cmdutil/factory"
	"github.com/4rchr4y/bpm/cli/cmdutil/require"
	"github.com/4rchr4y/bpm/constant"
	"github.com/spf13/cobra"
)

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
				Author: func() *bundlefile.AuthorExpr {
					user, err := f.GitCLI.User()
					if err != nil {
						return nil
					}

					return &bundlefile.AuthorExpr{
						Username: user.Username,
						Email:    user.Email,
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
	Author     *bundlefile.AuthorExpr                                 // git information about the author
	Encoder    *encode.Encoder                                        // decoder of bundle component files
	WriteFile  func(name string, data []byte, perm fs.FileMode) error // func of saving a file to disk
}

func initRun(opts *initOptions) error {
	bundlefileContent := bundleFileContent(opts.Encoder, opts.Repository, opts.Author)
	bundlefileHash := md5.Sum(bytes.TrimSpace(bundlefileContent))
	// TODO: to use manifest util to generate init lockfile data
	lockfileContent := lockfileContent(opts.Encoder, hex.EncodeToString(bundlefileHash[:]), "2024")

	files := map[string][]byte{
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
