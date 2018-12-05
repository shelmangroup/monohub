package storage

import (
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Storage struct {
	Root string
	Repo *GitRepository
}

var (
	dataDir = kingpin.Flag("data-directory", "Data directory").Short('d').Required().String()
)

func NewStorage() *Storage {
	log.WithField("data-dir", *dataDir).Info("Initializing storage")
	r := ensureDirectory(*dataDir)

	storage := &Storage{
		Root: r,
	}

	storage.Repo = NewGitRepository(path.Join(r, "repo"))
	err := storage.Repo.Init()
	if err != nil {
		log.WithError(err).Fatal("Error initializing git repo")
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
