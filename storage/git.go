package storage

import (
	"path/filepath"
)

type GitRepository struct {
	Path string
}

func NewGitRepository(path string) *GitRepository {
	p := filepath.Clean(path)

	repo := &GitRepository{
		Path: p,
	}

	return repo
}
