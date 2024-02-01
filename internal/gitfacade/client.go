package gitfacade

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

type GitFacade struct{}

func NewGitFacade() *GitFacade {
	return &GitFacade{}
}

func (gf *GitFacade) CloneWithContext(ctx context.Context, opts *git.CloneOptions) (*git.Repository, error) {
	repo, err := git.CloneContext(ctx, memory.NewStorage(), nil, opts)
	if err != nil {
		return nil, err
	}

	return repo, err
}

type User struct {
	Username string
	Email    string
}

func (gf *GitFacade) GetUser() (user User, err error) {
	user.Username, err = gf.getGitUserInfo("username")
	if err != nil {
		return user, err
	}

	user.Email, err = gf.getGitUserInfo("email")
	if err != nil {
		return user, err
	}

	return user, nil
}

func (gf *GitFacade) getGitUserInfo(key string) (string, error) {
	// TODO: use exec as interface
	cmd := exec.Command("git", "config", "--get", fmt.Sprintf("user.%s", key))
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
