package server

import (
	"context"

	pb "github.com/shelmangroup/monohub/api"
	log "github.com/sirupsen/logrus"
	// "gopkg.in/src-d/go-git.v4"
	// "gopkg.in/src-d/go-git.v4/plumbing"
	// "gopkg.in/src-d/go-git.v4/plumbing/object"
)

// Commits return commit information
func (s *Server) Blobs(ctx context.Context, req *pb.BlobRequest) (*pb.BlobResponse, error) {
	log.WithField("context", ctx).Debug("Blobs called")
	return nil, nil
}
