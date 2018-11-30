package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4/plumbing/protocol/packp"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/server"
)

func StartServer() error {
	log.Info("Starting server on :8822")
	wd, _ := os.Getwd()

	return http.ListenAndServe(":8822", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")

		// File Backend
		s := strings.SplitN(strings.TrimLeft(r.URL.Path, "/"), "/", 2)
		ep, err := transport.NewEndpoint(path.Join("file://", wd, s[0]))
		if err != nil {
			log.Fatalf("Tansport error: %v", err)
		}
		ups, err := server.DefaultServer.NewUploadPackSession(ep, nil)
		if err != nil {
			log.Fatalf("NewServer error: %v", err)
		}

		// GCS Backend
		// ep, err := transport.NewEndpoint("http://")
		// if err != nil {
		// 	log.Fatalf("Tansport error: %v", err)
		// }
		// ups, err := server.NewServer(NewGCSStorageLoader()).NewUploadPackSession(ep, nil)
		// if err != nil {
		// 	log.Fatalf("NewServer error: %v", err)
		// }

		if strings.Contains(r.URL.Path, "info") {
			advs, err := ups.AdvertisedReferences()
			if err != nil {
				log.Fatalf("AdvertisedRef error: %v", err)
			}
			advs.Prefix = [][]byte{
				[]byte("# service=git-upload-pack"),
				[]byte(""),
			}
			w.Header().Set("Content-Type", "application/x-git-upload-pack-advertisement")
			log.Debugf("%+v", w)
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

		up, err := ups.UploadPack(r.Context(), upakreq)
		if err != nil {
			log.Fatalf("UploadPack error: %v", err)
		}
		w.Header().Set("Content-Type", "application/x-git-upload-pack-result")
		up.Encode(w)
	}))
}
