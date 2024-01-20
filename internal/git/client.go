package git

import (
	"net/url"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
)

type GitClient struct{}

func NewClient() *GitClient {
	return &GitClient{}
}

type PlainCloneInput struct {
	Dir     string
	IsBare  bool
	Options *git.CloneOptions
}

func (gs *GitClient) PlainClone(input *PlainCloneInput) (*git.Repository, error) {
	repo, err := git.PlainClone(input.Dir, input.IsBare, input.Options)
	if err != nil {
		return nil, err
	}

	return repo, err
}

func extractRepoName(repoURL string) (string, error) {
	parsedURL, err := url.Parse(repoURL)
	if err != nil {
		return "", err
	}

	repoName := path.Base(parsedURL.Path)
	return strings.TrimSuffix(repoName, ".git"), nil
}
