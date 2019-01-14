package server

import (
	"context"

	pb "github.com/shelmangroup/monohub/api"
	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// Trees return tree information
func (s *Server) Trees(ctx context.Context, req *pb.TreeRequest) (*pb.TreeResponse, error) {
	log.WithField("context", ctx).Debug("Blobs called")
	r, err := git.PlainOpen(s.storage.Repo.Path)
	if err != nil {
		return nil, err
	}

	tree, err := r.TreeObject(plumbing.NewHash(req.Sha))
	if err != nil {
		return nil, err
	}
	var files []*pb.TreeFile
	err = tree.Files().ForEach(func(f *object.File) error {
		file := &pb.TreeFile{
			Path: f.Name,
			Mode: f.Mode.String(),
			Sha:  f.Hash.String(),
			Size: f.Size,
			Type: f.Type().String(),
		}
		files = append(files, file)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &pb.TreeResponse{
		Sha:       tree.Hash.String(),
		Tree:      files,
		Truncated: false,
	}, nil
}
