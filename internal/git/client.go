package git

import (
	"context"
	"net/url"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

type GitClient struct{}

func NewClient() *GitClient {
	return &GitClient{}
}

func (gs *GitClient) CloneWithContext(ctx context.Context, opts *git.CloneOptions) (*git.Repository, error) {
	repo, err := git.CloneContext(ctx, memory.NewStorage(), nil, opts)
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
