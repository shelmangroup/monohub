package server

import (
	"net/http"

	"github.com/shelmangroup/monohub/storage"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Server struct {
	storage    *storage.Storage
	gitHandler *GitHandler
	httpServer *http.Server
}

var (
	listenAddress = kingpin.Flag("listen-address", "HTTP address").Default(":8822").String()
)

func NewServer(storage *storage.Storage) *Server {

	server := &Server{
		storage: storage,
	}
	server.gitHandler = NewGitHandler(storage)

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
