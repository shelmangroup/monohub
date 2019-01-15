package server

import (
	"context"

	pb "github.com/shelmangroup/monohub/api"
	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// Tags return tag information
func (s *Server) Tags(ctx context.Context, req *pb.TagRequest) (*pb.TagResponse, error) {
	log.WithField("context", ctx).Debug("Tags called")
	r, err := git.PlainOpen(s.storage.Repo.Path)
	if err != nil {
		return nil, err
	}

	tag, err := r.TagObject(plumbing.NewHash(req.Sha))
	if err != nil {
		return nil, err
	}

	obj := &pb.TagObject{
		Type: tag.Type().String(),
		Sha:  tag.Hash.String(),
	}

	return &pb.TagResponse{
		Tag:    tag.Name,
		Sha:    tag.Hash.String(),
		Object: obj,
		Tagger: &pb.Tagger{
			Name:  tag.Tagger.Name,
			Email: tag.Tagger.Email,
		},
		Message: tag.Message,
	}, nil
}
