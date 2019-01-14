package server

import (
	"context"

	pb "github.com/shelmangroup/monohub/api"
	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// References return refs information
func (s *Server) References(ctx context.Context, req *pb.ReferenceRequest) (*pb.ReferenceResponse, error) {
	log.WithField("context", ctx).Debug("References called")
	r, err := git.PlainOpen(s.storage.Repo.Path)
	if err != nil {
		return nil, err
	}

	ref, err := r.Reference(plumbing.NewBranchReferenceName(req.Ref), true)
	if err != nil {
		return nil, err
	}
	var objects []*pb.RefObject
	obj := &pb.RefObject{
		Type: ref.Type().String(),
		Sha:  ref.Hash().String(),
	}
	objects = append(objects, obj)

	return &pb.ReferenceResponse{
		Ref:     ref.Hash().String(),
		Objects: objects,
	}, nil
}
