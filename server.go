package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4/plumbing/protocol/packp"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/server"
)

func StartServer() error {
	log.Info("Starting server on :8822")
	return http.ListenAndServe(":8822", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")

		ep, _ := transport.NewEndpoint("http://")
		ups, _ := server.NewServer(NewGCSStorageLoader()).NewUploadPackSession(ep, nil)
		if strings.Contains(r.URL.Path, "info") {
			advs, _ := ups.AdvertisedReferences()
			advs.Prefix = [][]byte{
				[]byte("# service=git-upload-pack"),
				[]byte(""),
			}
			w.Header().Set("Content-Type", "application/x-git-upload-pack-advertisement")
			advs.Encode(w)
			return
		}
		defer r.Body.Close()
		var rdr io.ReadCloser = r.Body

		if r.Header.Get("Content-Encoding") == "gzip" {
			rdr, _ = gzip.NewReader(r.Body)
		}

		upakreq := packp.NewUploadPackRequest()
		upakreq.Decode(rdr)

		up, _ := ups.UploadPack(r.Context(), upakreq)
		w.Header().Set("Content-Type", "application/x-git-upload-pack-result")
		up.Encode(w)
	}))
}
