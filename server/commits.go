package server

import (
	"context"

	"github.com/golang/protobuf/ptypes"
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
		log.WithField("context", ctx).Errorf("Commit: %s not found", req.Sha)
		return nil, err
	}

	currentTree, err := c.Tree()
	if err != nil {
		return nil, err
	}

	pIter := c.Parents()
	var parents []*pb.Parent
	var files []*pb.File
	err = pIter.ForEach(func(p *object.Commit) error {
		parent := &pb.Parent{
			Sha: p.Hash.String(),
		}
		parents = append(parents, parent)

		prevTree, err := p.Tree()
		if err != nil {
			return err
		}

		changes, err := currentTree.Diff(prevTree)
		if err != nil {
			return err
		}

		for _, ch := range changes {
			p, err := ch.Patch()
			if err != nil {
				return err
			}
			patch := p.String()

			var addition int
			var deletion int
			for _, fs := range p.Stats() {
				addition += fs.Addition
				deletion += fs.Deletion
			}

			_, to, err := ch.Files()
			if err != nil {
				return err
			}

			file := &pb.File{
				Filename:  ch.To.Name,
				Additions: int64(addition),
				Deletions: int64(deletion),
				Patch:     patch,
				BlobUrl:   to.Hash.String(),
			}
			files = append(files, file)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	commitTimestamp, err := ptypes.TimestampProto(c.Committer.When)
	if err != nil {
		return nil, err
	}

	return &pb.CommitResponse{
		Sha: c.Hash.String(),
		Author: &pb.Author{
			Id:    1,
			Login: c.Author.Email,
			Name:  c.Author.Name,
			Email: c.Author.Email,
		},
		Commit: &pb.Commit{
			Committer: &pb.Author{
				Id:    1,
				Login: c.Committer.Email,
				Name:  c.Committer.Name,
				Email: c.Committer.Email,
				Date:  commitTimestamp,
			},
			Message: c.Message,
			Tree: &pb.Tree{
				Sha: c.TreeHash.String(),
			},
		},
		Files:   files,
		Stats:   &pb.Stats{},
		Parents: parents,
	}, nil
}
