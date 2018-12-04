package server

import (
	"bytes"
	"net/http"
	"net/http/cgi"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type GitHandler struct {
	CGI    string
	server *Server
}

func NewGitHandler(server *Server) *GitHandler {
	dir, err := exec.Command("git", "--exec-path").Output()
	if err != nil {
		log.Fatal("Could not find git binary")
	}
	dir = bytes.TrimRight(dir, "\r\n")
	cgi := filepath.Join(string(dir), "git-http-backend")
	log.Debugf("Git CGI command line: %s", cgi)

	return &GitHandler{
		CGI:    cgi,
		server: server,
	}
}

func (h *GitHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {

	env := []string{
		"GIT_PROJECT_ROOT=" + h.server.RepoPath,
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
		Path:   h.CGI,
		Root:   "/",
		Env:    env,
		Stderr: &stdErr,
	}
	handler.ServeHTTP(w, r)

	if stdErr.Len() > 0 {
		log.Infof("[git] %s", stdErr.String())
	}
}
