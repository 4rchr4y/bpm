package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/4rchr4y/bpm/bundle"
	"github.com/4rchr4y/bpm/cli/require"
	"github.com/4rchr4y/bpm/internal/encode"
	"github.com/4rchr4y/bpm/manager"
	"github.com/4rchr4y/godevkit/syswrap"
	"github.com/spf13/cobra"
)

const initCmdDesc = `
The 'bpm init' command is designed to generate a new bundle directory 
structure, complete with standard files typical for a bundle.

For instance, executing bpm new foo would result in a directory 
setup similar to the following:

foo/
├── .bpmignore		# Ignore files for bpm.
├── .gitignore		# Ignore for git system.
└── bundle.toml		# File with bundle information.

'bpm init' takes a path for an argument. If directories in the given path
do not exist, bpm will attempt to create them as it goes. If the given
destination exists and there are files in that directory, conflicting files
will be overwritten, but other files will be left alone.
`

func newInitCmd(args []string) *cobra.Command {
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
		Short: "Init a new bundle with.",
		Long:  initCmdDesc,
		RunE:  runInitCmd,
	}

	return cmd
}

func runInitCmd(cmd *cobra.Command, args []string) error {
	osWrap := new(syswrap.OsWrapper)
	tomlEncoder := encode.NewTomlEncoder()

	bpmManager := manager.NewBpm()
	initCmd := manager.NewInitCommand(&manager.InitCmdResources{
		OsWrap:      osWrap,
		TomlEncoder: tomlEncoder,
	})

	if err := bpmManager.RegisterCommand(
		initCmd,
	); err != nil {
		return err
	}

	if _, err := manager.ExecuteInitCmd(initCmd, &manager.InitCmdInput{
		Name:   args[0],
		Author: buildAuthorInfo(),
	}); err != nil {
		return err
	}

	return nil
}

func buildAuthorInfo() *bundle.AuthorExpr {
	username, _ := getGitUserInfo("username")
	email, _ := getGitUserInfo("email")

	if username == "" || email == "" {
		return nil
	}

	return &bundle.AuthorExpr{
		Username: username,
		Email:    email,
	}
}

func getGitUserInfo(target string) (string, error) {
	cmd := exec.Command("git", "config", "--get", fmt.Sprintf("user.%s", target))
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
