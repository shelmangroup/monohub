package server

import (
	"context"
	"encoding/base64"
	"io/ioutil"

	pb "github.com/shelmangroup/monohub/api"
	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	// "gopkg.in/src-d/go-git.v4/plumbing/object"
)

// Blobs return commit information
func (s *Server) Blobs(ctx context.Context, req *pb.BlobRequest) (*pb.BlobResponse, error) {
	log.WithField("context", ctx).Debug("Blobs called")
	r, err := git.PlainOpen(s.storage.Repo.Path)
	if err != nil {
		return nil, err
	}
	blob, err := r.BlobObject(plumbing.NewHash(req.Sha))
	if err != nil {
		return nil, err
	}
	br, err := blob.Reader()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(br)
	if err != nil {
		return nil, err
	}

	return &pb.BlobResponse{
		Content:  base64.StdEncoding.EncodeToString(data),
		Encoding: "base64",
		Size:     blob.Size,
		Sha:      blob.Hash.String(),
	}, nil
}
