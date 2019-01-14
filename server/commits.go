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
		log.WithField("context", ctx).Errorf("Commit: %s not found", req.Sha)
		return nil, err
	}
	prevCommit, err := c.Parents().Next()
	if err != nil {
		return nil, err
	}
	prevTree, err := prevCommit.Tree()
	if err != nil {
		return nil, err
	}
	currentTree, err := c.Tree()
	if err != nil {
		return nil, err
	}

	changes, err := currentTree.Diff(prevTree)
	if err != nil {
		return nil, err
	}

	var files []*pb.File
	var p *object.Patch
	var patch string
	for _, ch := range changes {
		p, err = ch.Patch()
		if err != nil {
			return nil, err
		}
		patch = p.String()

		var addition int
		var deletion int
		for _, fs := range p.Stats() {
			addition = fs.Addition
			deletion = fs.Deletion
		}

		_, to, err := ch.Files()
		if err != nil {
			return nil, err
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
