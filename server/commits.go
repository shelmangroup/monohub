package server

import (
	"context"

	pb "github.com/shelmangroup/monohub/api"
	log "github.com/sirupsen/logrus"
)

func (s *Server) Commits(ctx context.Context, req *pb.CommitRequest) (*pb.CommitResponse, error) {
	log.WithField("context", ctx).Debug("Commits called")
	return nil, nil
}
