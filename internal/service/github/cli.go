package github

import (
	"fmt"
	"os/exec"
	"strings"
)

type GitCLI struct{}

type User struct {
	Username string
	Email    string
}

func (cli *GitCLI) User() (user User, err error) {
	user.Username, err = getGitUserInfo("username")
	if err != nil {
		return user, err
	}

	user.Email, err = getGitUserInfo("email")
	if err != nil {
		return user, err
	}

	return user, nil
}

func getGitUserInfo(key string) (string, error) {
	// TODO: use exec as interface
	cmd := exec.Command("git", "config", "--get", fmt.Sprintf("user.%s", key))
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
