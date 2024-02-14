package github

import (
	"context"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

type GitClient struct{}

func (client *GitClient) CloneWithContext(ctx context.Context, opts *git.CloneOptions) (*git.Repository, error) {
	repo, err := git.CloneContext(ctx, memory.NewStorage(), nil, opts)
	if err != nil {
		return nil, err
	}

	return repo, err
}
