package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
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

func (r *GitRepository) SetupHooks(grpcSocket string) error {
	log.Info("Initializing git hooks")
	hooksDir := ensureDirectory(path.Join(r.Path, "hooks"))
	log.WithField("hooksDir", hooksDir).Debug("git hook directory initialized")

	r.setupHook(hooksDir, "pre-receive", grpcSocket)
	r.setupHook(hooksDir, "post-receive", grpcSocket)
	return nil
}

func (r *GitRepository) setupHook(hooksDir string, hook string, grpcSocket string) {
	binPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.WithError(err).Fatal("Error getting path for hooks binary")
	}

	data := []byte(fmt.Sprintf("#!/bin/sh\nexec %s hook %s --log-level=%s -g %s\n", binPath, hook, log.GetLevel(), grpcSocket))
	err = ioutil.WriteFile(path.Join(hooksDir, hook), data, 0755)
	if err != nil {
		log.WithError(err).Fatalf("Error initializing git hook %s", hook)
	}
}

func (r *GitRepository) Init() error {
	log.WithField("path", r.Path).Info("Initializing git repository")
	_, err := git.PlainInit(r.Path, true)
	if err == git.ErrRepositoryAlreadyExists {
		return nil
	}
	return err
}

func (r *GitRepository) GetGitRepo() *git.Repository {
	repo, err := git.PlainOpen(r.Path)
	if err != nil {
		log.WithError(err).Fatal("Error opening git repository")
	}
	return repo
}
