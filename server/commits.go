package server

import (
	"context"

	pb "github.com/shelmangroup/monohub/api"
	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// Commits return commit information
func (s *Server) Commits(ctx context.Context, req *pb.CommitRequest) (*pb.CommitResponse, error) {
	log.WithField("context", ctx).Debug("Commits called")
	r, err := git.PlainOpen(s.storage.Repo.Path)
	if err != nil {
		return nil, err
	}

	c, err := r.CommitObject(plumbing.NewHash(req.Sha))
	if err != nil {
		return nil, err
	}

	tree, err := c.Tree()
	if err != nil {
		return nil, err
	}

	tIter := tree.Files()
	if err != nil {
		return nil, err
	}

	var files []*pb.File
	err = tIter.ForEach(func(f *object.File) error {
		file := &pb.File{
			Filename: f.Name,
		}
		files = append(files, file)
		return nil
	})
	if err != nil {
		return nil, err
	}

	pIter := c.Parents()
	var parents []*pb.Parent
	err = pIter.ForEach(func(c *object.Commit) error {
		parent := &pb.Parent{
			Sha: c.Hash.String(),
		}
		parents = append(parents, parent)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &pb.CommitResponse{
		Sha: c.Hash.String(),
		Author: &pb.Author{
			Id:    1,
			Login: c.Author.Name,
			Name:  c.Author.Name,
			Email: c.Author.Email,
		},
		Commit: &pb.Commit{
			Committer: &pb.Author{
				Id:    1,
				Login: c.Author.Name,
				Name:  c.Author.Name,
				Email: c.Author.Email,
			},
			Message: c.Message,
			Tree: &pb.Tree{
				Sha: tree.Hash.String(),
			},
		},
		Files:   files,
		Stats:   &pb.Stats{},
		Parents: parents,
	}, nil
}
