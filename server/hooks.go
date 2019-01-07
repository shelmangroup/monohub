package server

import (
	"context"
	"fmt"
	"net"
	"os"

	api "github.com/shelmangroup/monohub/api"
	"github.com/shelmangroup/monohub/util"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func (s *Server) RunHooksServer() error {
	sock := s.storage.HooksSocketPath()
	log.WithField("socket", sock).Info("Starting git hooks grpc server")

	err := os.Remove(sock)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	listen, err := net.Listen("unix", sock)
	if err != nil {
		return err
	}
	gs := grpc.NewServer()
	api.RegisterGitHooksServer(gs, s)
	reflection.Register(gs)
	err = gs.Serve(listen)
	return err
}

func (s *Server) PreReceive(ctx context.Context, req *api.PreReceiveRequest) (*api.HookResult, error) {
	log.WithField("req", req).Debug("PreReceive called")
	reply := &api.HookResult{
		Status:  api.HookStatus_OK,
		Message: "ok",
	}

	log.Debug("Disallowing force-push on master branch")

	repo := s.storage.Repo.GetGitRepo()

	for _, op := range req.GetOps() {
		log.Debugf("Validating operation: %s", op)
		if op.OldValue == "0000000000000000000000000000000000000000" {
			log.Debug("No commit history. Allowing operation.")
			continue
		}

		commit, err := repo.CommitObject(plumbing.NewHash(op.NewValue))
		if err != nil {
			msg := fmt.Sprintf("Error resolving revision: %s", op.NewValue)
			log.WithError(err).Errorf(msg)
			return &api.HookResult{Status: api.HookStatus_ERROR, Message: msg}, err
		}
		log.Debugf("Commit: %s", commit)

		hashes, err := util.MergeBase(repo, plumbing.NewHash(op.NewValue), plumbing.NewHash(op.OldValue))
		if len(hashes) == 0 {
			log.Warn("Nope, no common commit!")
			reply = &api.HookResult{Status: api.HookStatus_ERROR, Message: "Force push is a no-no"}
		} else {
			log.Infof("Alrighty! Common ancestors: %s", hashes)
			reply = &api.HookResult{
				Status:  api.HookStatus_OK,
				Message: "ok",
			}
			break
		}

	}

	return reply, nil
}
