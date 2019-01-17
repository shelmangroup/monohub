package lfs

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const (
	contentMediaType = "application/vnd.git-lfs"
	metaMediaType    = contentMediaType + "+json"
)

// Server links ContentStore, and MetaStore to provide the LFS server.
type Server struct {
	contentStore *ContentStore
	metaStore    *MetaStore
}

// NewServer creates a new App using the ContentStore and MetaStore provided
func NewServer(content *ContentStore, meta *MetaStore) *Server {
	return &Server{
		contentStore: content,
		metaStore:    meta,
	}
}

func (s *Server) LocksHandler(w http.ResponseWriter, r *http.Request) {
	return
}
func (s *Server) LocksVerifyHandler(w http.ResponseWriter, r *http.Request) {
	return
}
func (s *Server) CreateLockHandler(w http.ResponseWriter, r *http.Request) {
	return
}
func (s *Server) DeleteLockHandler(w http.ResponseWriter, r *http.Request) {
	return
}
func (s *Server) BatchHandler(w http.ResponseWriter, r *http.Request) {
	return
}
func (s *Server) GetContentHandler(w http.ResponseWriter, r *http.Request) {
	return
}
func (s *Server) GetMetaHandler(w http.ResponseWriter, r *http.Request) {
	return
}
func (s *Server) PutHandler(w http.ResponseWriter, r *http.Request) {
	return
}
func (s *Server) PostHandler(w http.ResponseWriter, r *http.Request) {
	return
}
func (s *Server) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	return
}

// ContentMatcher provides a mux.MatcherFunc that only allows requests that
// contain an Accept header with the contentMediaType
func ContentMatcher(r *http.Request, m *mux.RouteMatch) bool {
	mediaParts := strings.Split(r.Header.Get("Accept"), ";")
	mt := mediaParts[0]
	return mt == contentMediaType
}

// MetaMatcher provides a mux.MatcherFunc that only allows requests
// that contain an Accept header with the metaMediaType
func MetaMatcher(r *http.Request, m *mux.RouteMatch) bool {
	mediaParts := strings.Split(r.Header.Get("Accept"), ";")
	mt := mediaParts[0]
	return mt == metaMediaType
}
