package main

import (
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/server"
)

type gcsLoader struct{}

func NewGCSStorageLoader() server.Loader {
	return &gcsLoader{}
}

func (l *gcsLoader) Load(ep *transport.Endpoint) (storer.Storer, error) {
	s, err := NewStorage()
	if err != nil {
		return nil, err
	}
	return s, nil
}
