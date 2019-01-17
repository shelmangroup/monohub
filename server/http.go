package server

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	api "github.com/shelmangroup/monohub/api"
	"github.com/shelmangroup/monohub/server/lfs"
	"github.com/shelmangroup/monohub/storage"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"go.opencensus.io/stats/view"
	"google.golang.org/grpc"
)

type HttpServer struct {
	storage *storage.Storage
	server  *Server
	lfs     *lfs.Server
}

func NewHttpServer(server *Server) *HttpServer {
	return &HttpServer{
		server:  server,
		storage: server.storage,
	}
}

func (s *HttpServer) Run() error {
	if err := view.Register(ochttp.DefaultServerViews...); err != nil {
		return err
	}

	gwmux := runtime.NewServeMux()
	ctx := context.Background()
	dopts := []grpc.DialOption{grpc.WithInsecure()}
	err := api.RegisterMonoHubHandlerFromEndpoint(ctx, gwmux, *grpcAddress, dopts)
	if err != nil {
		log.WithError(err).Error("grpc gateway register failed")
		return err
	}

	router := mux.NewRouter()
	router.HandleFunc("/info/refs", s.getInfoRefsHandler).Queries("service", "{service}").Methods("GET")
	router.HandleFunc("/git-upload-pack", s.uploadPackHandler).Methods("POST")
	router.HandleFunc("/git-receive-pack", s.receivePackHandler).Methods("POST")

	// GIT LFS
	router.HandleFunc("/{repo}/locks", s.lfs.LocksHandler).Methods("GET").MatcherFunc(lfs.MetaMatcher)
	router.HandleFunc("/{repo}/locks/verify", s.lfs.LocksVerifyHandler).Methods("POST").MatcherFunc(lfs.MetaMatcher)
	router.HandleFunc("/{repo}/locks", s.lfs.CreateLockHandler).Methods("POST").MatcherFunc(lfs.MetaMatcher)
	router.HandleFunc("/{repo}/locks/{id}/unlock", s.lfs.DeleteLockHandler).Methods("POST").MatcherFunc(lfs.MetaMatcher)

	router.HandleFunc("/objects/{oid}", s.lfs.GetContentHandler).Methods("GET", "HEAD").MatcherFunc(lfs.ContentMatcher)
	router.HandleFunc("/objects/{oid}", s.lfs.GetMetaHandler).Methods("GET", "HEAD").MatcherFunc(lfs.MetaMatcher)
	router.HandleFunc("/objects/{oid}", s.lfs.PutHandler).Methods("PUT").MatcherFunc(lfs.ContentMatcher)
	router.HandleFunc("/objects", s.lfs.PostHandler).Methods("POST").MatcherFunc(lfs.MetaMatcher)
	router.HandleFunc("/verify/{oid}", s.lfs.VerifyHandler).Methods("POST")
	// END GIT LFS

	router.PathPrefix("/").Handler(gwmux)

	log.WithField("address", *listenAddress).Info("Starting HTTP server")
	return http.ListenAndServe(*listenAddress,
		&ochttp.Handler{
			Handler:     router,
			Propagation: &b3.HTTPFormat{},
		},
	)
}

func (s *HttpServer) setNoCacheHeaders(w http.ResponseWriter) {
	w.Header().Set("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-cache, max-age=0, must-revalidate")
}

func (s *HttpServer) getInfoRefsHandler(w http.ResponseWriter, req *http.Request) {
	log.Debug("infoRefs")
	s.setNoCacheHeaders(w)

	service := getServiceType(req)
	cmd := exec.Command("git", service, "--stateless-rpc", "--advertise-refs", s.storage.Repo.Path)
	cmd.Dir = s.storage.Repo.Path

	refs, err := cmd.Output()
	if err != nil {
		log.Errorf("fail to serve RPC(%s): %v", service, err)
		return
	}

	w.Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-advertisement", service))
	w.WriteHeader(http.StatusOK)
	w.Write(packetWrite("# service=git-" + service + "\n"))
	w.Write([]byte("0000"))
	w.Write(refs)
}

func (s *HttpServer) uploadPackHandler(w http.ResponseWriter, req *http.Request) {
	log.Debug("uploadPack")
	s.serviceRPC("upload-pack", w, req)
}

func (s *HttpServer) receivePackHandler(w http.ResponseWriter, req *http.Request) {
	log.Debug("receivePack")
	s.serviceRPC("receive-pack", w, req)
}

func (s *HttpServer) serviceRPC(service string, w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	w.Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-result", service))

	var err error
	var reqBody = req.Body

	// Handle GZIP.
	if req.Header.Get("Content-Encoding") == "gzip" {
		reqBody, err = gzip.NewReader(reqBody)
		if err != nil {
			log.Errorf("fail to create gzip reader: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	var stderr bytes.Buffer
	cmd := exec.Command("git", service, "--stateless-rpc", s.storage.Repo.Path)
	cmd.Dir = s.storage.Repo.Path

	cmd.Stdout = w
	cmd.Stdin = reqBody
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("fail to serve RPC(%s): %v - %s", service, err, stderr.String())
		return
	}
	return
}

func packetWrite(str string) []byte {
	s := strconv.FormatInt(int64(len(str)+4), 16)
	if len(s)%4 != 0 {
		s = strings.Repeat("0", 4-len(s)%4) + s
	}
	return []byte(s + str)
}

func getServiceType(r *http.Request) string {
	serviceType := r.FormValue("service")
	if !strings.HasPrefix(serviceType, "git-") {
		return ""
	}
	return strings.Replace(serviceType, "git-", "", 1)
}
