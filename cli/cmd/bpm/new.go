package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/4rchr4y/bpm/cli/require"
	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/godevkit/syswrap"
	"github.com/spf13/cobra"
)

const newCmdDesc = `
The 'bpm new' command is designed to generate a new bundle directory 
structure, complete with standard files typical for a bundle.

For instance, executing bpm new foo would result in a directory 
setup similar to the following:

foo/
├── .gitignore		# Ignore for git system.
└── bundle.toml		# File with bundle information.

'bpm new' takes a path for an argument. If directories in the given path
do not exist, bpm will attempt to create them as it goes. If the given
destination exists and there are files in that directory, conflicting files
will be overwritten, but other files will be left alone.
`

func newNewCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "new NAME",
		Aliases: []string{"n", "create"},
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
		Short: "Create a new bundle with the given name",
		Long:  newCmdDesc,
		RunE:  runNewCmd,
	}

	return cmd
}

func runNewCmd(cmd *cobra.Command, args []string) error {
	osWrap := new(syswrap.OsWrapper)

	dirPath := args[0]
	if err := osWrap.MkdirAll(dirPath, 0755); err != nil {
		return err
	}

	files := map[string]string{
		".gitignore":            gitignoreFileContent(),
		constant.BundleFileName: bundleFileContent(filepath.Base(dirPath), buildAuthorInfo()),
	}

	for fileName, content := range files {
		fullPath := filepath.Join(dirPath, fileName)
		if err := osWrap.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file '%s': %v", fullPath, err)
		}
	}

	return nil
}

func buildAuthorInfo() string {
	username, _ := getGitUserInfo("username")
	email, _ := getGitUserInfo("email")

	if username == "" || email == "" {
		return ""
	}

	return fmt.Sprintf("%s %s", username, email)
}

func getGitUserInfo(target string) (string, error) {
	cmd := exec.Command("git", "config", "--get", fmt.Sprintf("user.%s", target))
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func bundleFileContent(name string, author string) string {
	return fmt.Sprintf(`[package]
name = "%s"
author = ["%s"]
description = ""
	
[dependencies]`, name, author,
	)
}

func gitignoreFileContent() string {
	return `bundle.lock`
}
