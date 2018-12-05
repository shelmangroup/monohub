package storage

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
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

func (r *GitRepository) Init() error {
	log.WithField("path", r.Path).Info("Initializing git repository")
	_, err := git.PlainInit(r.Path, true)

	if err == git.ErrRepositoryAlreadyExists {
		return nil
	}

	return err
}
