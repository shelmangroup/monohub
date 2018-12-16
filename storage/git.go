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

func (r *GitRepository) initHooks() error {
	log.Info("Initializing git hooks")
	hooksDir := ensureDirectory(path.Join(r.Path, "hooks"))
	log.WithField("hooksDir", hooksDir).Debug("git hook directory initialized")

	binPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return err
	}

	data := []byte(fmt.Sprintf("#!/bin/sh\nexec %s hook pre-receive -d %s\n", binPath, r.Path))
	err = ioutil.WriteFile(path.Join(hooksDir, "pre-receive"), data, 0755)
	if err != nil {
		return err
	}

	data = []byte(fmt.Sprintf("#!/bin/sh\nexec %s hook post-receive -d %s\n", binPath, r.Path))
	err = ioutil.WriteFile(path.Join(hooksDir, "post-receive"), data, 0755)
	if err != nil {
		return err
	}

	return nil
}

func (r *GitRepository) Init() error {
	log.WithField("path", r.Path).Info("Initializing git repository")
	_, err := git.PlainInit(r.Path, true)

	if err != nil && err != git.ErrRepositoryAlreadyExists {
		return err
	}

	err = r.initHooks()

	return err
}
