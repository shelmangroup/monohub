package server

import (
	"context"

	pb "github.com/shelmangroup/monohub/api"
	log "github.com/sirupsen/logrus"
)

func (s *Server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingReply, error) {
	log.WithField("context", ctx).Debug("Ping called")
	reply := &pb.PingReply{
		Version: "master",
	}
	return reply, nil
}
