package server

import (
	"bytes"
	"net/http"
	"net/http/cgi"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

type GitServer struct {
	CGI        string
	RepoPath   string
	httpServer *http.Server
}

var (
	listenAddress = kingpin.Flag("listen-address", "HTTP address").Default(":8822").String()
	repoPath      = kingpin.Flag("repository-path", "Git repository path").Default(".").String()
)

func NewGitServer() *GitServer {
	dir, err := exec.Command("git", "--exec-path").Output()
	if err != nil {
		log.Fatal("Could not find git binary")
	}
	dir = bytes.TrimRight(dir, "\r\n")
	cgi := filepath.Join(string(dir), "git-http-backend")
	log.Debugf("Git CGI command line: %s", cgi)

	server := &GitServer{
		CGI:      cgi,
		RepoPath: *repoPath,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", server.gitHandler)
	server.httpServer = &http.Server{
		Addr:    *listenAddress,
		Handler: mux,
	}

	return server
}

func (s *GitServer) Serve() error {
	log.WithField("address", *listenAddress).Info("Starting server")

	return s.httpServer.ListenAndServe()
}

func (s *GitServer) gitHandler(w http.ResponseWriter, r *http.Request) {

	env := []string{
		"GIT_PROJECT_ROOT=" + s.RepoPath,
		"GIT_HTTP_EXPORT_ALL=1",
		"GIT_TRACE=2",
		"REMOTE_USER=dln",
		"GIT_HTTP_MAX_REQUEST_BUFFER=4000M",
	}

	// env = append(env, "REMOTE_USER="+username)

	log.Debug("Running git handler")
	log.Debug("Request: ", r)
	var stdErr bytes.Buffer
	handler := &cgi.Handler{
		Path:   s.CGI,
		Root:   "/",
		Env:    env,
		Stderr: &stdErr,
	}
	handler.ServeHTTP(w, r)

	if stdErr.Len() > 0 {
		log.Infof("[git] %s", stdErr.String())
	}
}
