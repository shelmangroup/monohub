package server

import (
	"context"
	"net"
	"os"

	api "github.com/shelmangroup/monohub/api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

func (s *Server) PreReceive(ctx context.Context, req *api.PreReceiveRequest) (*api.PreReceiveReply, error) {
	log.WithField("context", ctx).Debug("PreReceive called")
	reply := &api.PreReceiveReply{}
	return reply, nil
}
