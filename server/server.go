package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Server struct {
	RepoPath   string
	gitHandler *GitHandler
	httpServer *http.Server
}

var (
	listenAddress = kingpin.Flag("listen-address", "HTTP address").Default(":8822").String()
	repoPath      = kingpin.Flag("repository-path", "Git repository path").Default(".").String()
)

func NewServer() *Server {

	server := &Server{
		RepoPath: *repoPath,
	}
	server.gitHandler = NewGitHandler(server)

	mux := http.NewServeMux()
	mux.HandleFunc("/", server.gitHandler.HandleRequest)
	server.httpServer = &http.Server{
		Addr:    *listenAddress,
		Handler: mux,
	}

	return server
}

func (s *Server) Serve() error {
	log.WithField("address", *listenAddress).Info("Starting server")

	return s.httpServer.ListenAndServe()
}
