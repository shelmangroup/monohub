package storage

import (
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type Storage struct {
	Root string
	Repo *GitRepository
}

func NewStorage(dir string) *Storage {
	log.WithField("data-dir", dir).Info("Initializing storage")
	r := ensureDirectory(dir)

	storage := &Storage{
		Root: r,
	}

	storage.Repo = NewGitRepository(path.Join(r, "repo"))
	err := storage.Repo.Init()
	if err != nil {
		log.WithError(err).Fatal("Error initializing git repo")
	}
	err = storage.Repo.SetupHooks(storage.HooksSocketPath())
	if err != nil {
		log.WithError(err).Fatal("Error setting up git hooks")
	}

	return storage
}

func ensureDirectory(path string) string {
	p := filepath.Clean(path)

	l := log.WithField("path", p)

	s, err := os.Stat(p)
	if os.IsNotExist(err) {
		l.Info("Creating directory")
		err = os.MkdirAll(p, os.ModePerm)
		if err != nil {
			l.WithError(err).Fatal("Error creating directory")
		}
	} else if !s.IsDir() {
		l.Fatal("Path exists and is not a directory")
	}

	return p
}

func (s *Storage) HooksSocketPath() string {
	return path.Join(s.Root, "hooks_rpc.sock")
}
