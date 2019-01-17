package lfs

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/context"
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
	enc := json.NewEncoder(w)
	ll := &LockList{}

	w.Header().Set("Content-Type", metaMediaType)

	locks, nextCursor, err := s.metaStore.FilteredLocks("locks",
		r.FormValue("path"),
		r.FormValue("cursor"),
		r.FormValue("limit"))

	if err != nil {
		ll.Message = err.Error()
	} else {
		ll.Locks = locks
		ll.NextCursor = nextCursor
	}

	enc.Encode(ll)
}

func (s *Server) LocksVerifyHandler(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "USER")

	dec := json.NewDecoder(r.Body)
	enc := json.NewEncoder(w)

	w.Header().Set("Content-Type", metaMediaType)

	reqBody := &VerifiableLockRequest{}
	if err := dec.Decode(reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(&VerifiableLockList{Message: err.Error()})
		return
	}

	// Limit is optional
	limit := reqBody.Limit
	if limit == 0 {
		limit = 100
	}

	ll := &VerifiableLockList{}
	locks, nextCursor, err := s.metaStore.FilteredLocks("locks", "",
		reqBody.Cursor,
		strconv.Itoa(limit))
	if err != nil {
		ll.Message = err.Error()
	} else {
		ll.NextCursor = nextCursor

		for _, l := range locks {
			if l.Owner.Name == user {
				ll.Ours = append(ll.Ours, l)
			} else {
				ll.Theirs = append(ll.Theirs, l)
			}
		}
	}

	enc.Encode(ll)
}

func (s *Server) CreateLockHandler(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "USER").(string)

	dec := json.NewDecoder(r.Body)
	enc := json.NewEncoder(w)

	w.Header().Set("Content-Type", metaMediaType)

	var lockRequest LockRequest
	if err := dec.Decode(&lockRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(&LockResponse{Message: err.Error()})
		return
	}

	locks, _, err := s.metaStore.FilteredLocks("locks", lockRequest.Path, "", "1")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(&LockResponse{Message: err.Error()})
		return
	}
	if len(locks) > 0 {
		w.WriteHeader(http.StatusConflict)
		enc.Encode(&LockResponse{Message: "lock already created"})
		return
	}

	lock := &Lock{
		Id:       randomLockID(),
		Path:     lockRequest.Path,
		Owner:    User{Name: user},
		LockedAt: time.Now(),
	}

	if err := s.metaStore.AddLocks("locks", *lock); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(&LockResponse{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	enc.Encode(&LockResponse{
		Lock: lock,
	})
}

func (s *Server) DeleteLockHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lockID := vars["id"]
	user := context.Get(r, "USER").(string)

	dec := json.NewDecoder(r.Body)
	enc := json.NewEncoder(w)

	w.Header().Set("Content-Type", metaMediaType)

	var unlockRequest UnlockRequest

	if len(lockID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(&UnlockResponse{Message: "invalid lock id"})
		return
	}

	if err := dec.Decode(&unlockRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(&UnlockResponse{Message: err.Error()})
		return
	}

	l, err := s.metaStore.DeleteLock("locks", user, lockID, unlockRequest.Force)
	if err != nil {
		if err == errNotOwner {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		enc.Encode(&UnlockResponse{Message: err.Error()})
		return
	}
	if l == nil {
		w.WriteHeader(http.StatusNotFound)
		enc.Encode(&UnlockResponse{Message: "unable to find lock"})
		return
	}

	enc.Encode(&UnlockResponse{Lock: l})
}

func (s *Server) BatchHandler(w http.ResponseWriter, r *http.Request) {
	bv := unpackBatch(r)

	var responseObjects []*Representation

	url := fmt.Sprintf("%s://%s", r.URL.Scheme, r.URL.Host)

	// Create a response object
	done := make(chan *Representation, len(bv.Objects))

	for _, object := range bv.Objects {
		go func(object *RequestVars) {
			meta, err := s.metaStore.Get(object)
			if err == nil && s.contentStore.Exists(meta) { // Object is found and exists
				done <- s.Represent(object, meta, url, true, false)
				return
			}

			// Object is not found
			meta, err = s.metaStore.Put(object)
			if err == nil {
				done <- s.Represent(object, meta, url, meta.Existing, true)
				return
			}
			done <- nil
		}(object)
	}

	for i := 0; i < len(bv.Objects); i++ {
		responseObjects = append(responseObjects, <-done)
	}

	w.Header().Set("Content-Type", metaMediaType)

	respobj := &BatchResponse{Objects: responseObjects}
	json.NewEncoder(w).Encode(respobj)
}

func (s *Server) GetContentHandler(w http.ResponseWriter, r *http.Request) {
	rv := unpack(r)
	meta, err := s.metaStore.Get(rv)
	if err != nil {
		writeStatus(w, r, 404)
		return
	}

	// Support resume download using Range header
	var fromByte int64
	statusCode := 200
	if rangeHdr := r.Header.Get("Range"); rangeHdr != "" {
		regex := regexp.MustCompile(`bytes=(\d+)\-.*`)
		match := regex.FindStringSubmatch(rangeHdr)
		if match != nil && len(match) > 1 {
			statusCode = 206
			fromByte, _ = strconv.ParseInt(match[1], 10, 32)
			w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", fromByte, meta.Size-1, int64(meta.Size)-fromByte))
		}
	}

	content, err := s.contentStore.Get(meta, fromByte)
	if err != nil {
		writeStatus(w, r, 404)
		return
	}
	defer content.Close()

	w.WriteHeader(statusCode)
	io.Copy(w, content)
}

func (s *Server) GetMetaHandler(w http.ResponseWriter, r *http.Request) {
	rv := unpack(r)
	meta, err := s.metaStore.Get(rv)
	if err != nil {
		writeStatus(w, r, 404)
		return
	}

	w.Header().Set("Content-Type", metaMediaType)

	if r.Method == "GET" {
		url := fmt.Sprintf("%s://%s", r.URL.Scheme, r.URL.Host)
		enc := json.NewEncoder(w)
		enc.Encode(s.Represent(rv, meta, url, true, false))
	}
}

func (s *Server) PutHandler(w http.ResponseWriter, r *http.Request) {
	rv := unpack(r)
	meta, err := s.metaStore.Get(rv)
	if err != nil {
		writeStatus(w, r, 404)
		return
	}

	if err := s.contentStore.Put(meta, r.Body); err != nil {
		s.metaStore.Delete(rv)
		w.WriteHeader(500)
		fmt.Fprintf(w, `{"message":"%s"}`, err)
		return
	}
}

func (s *Server) PostHandler(w http.ResponseWriter, r *http.Request) {
	rv := unpack(r)
	meta, err := s.metaStore.Put(rv)
	if err != nil {
		writeStatus(w, r, 404)
		return
	}

	w.Header().Set("Content-Type", metaMediaType)

	sentStatus := 202
	if meta.Existing && s.contentStore.Exists(meta) {
		sentStatus = 200
	}
	w.WriteHeader(sentStatus)

	url := fmt.Sprintf("%s://%s", r.URL.Scheme, r.URL.Host)
	enc := json.NewEncoder(w)
	enc.Encode(s.Represent(rv, meta, url, meta.Existing, true))
}

func (s *Server) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	return
}

func (s *Server) Represent(rv *RequestVars, meta *MetaObject, url string, download, upload bool) *Representation {
	rep := &Representation{
		Oid:     meta.Oid,
		Size:    meta.Size,
		Actions: make(map[string]*link),
	}

	header := make(map[string]string)
	header["Accept"] = contentMediaType

	if len(rv.Authorization) > 0 {
		header["Authorization"] = rv.Authorization
	}

	if download {
		rep.Actions["download"] = &link{Href: url + "/objects/" + rep.Oid, Header: header}
	}

	if upload {
		rep.Actions["upload"] = &link{Href: url + "/objects/" + rep.Oid, Header: header}
	}
	return rep
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

func randomLockID() string {
	var id [20]byte
	rand.Read(id[:])
	return fmt.Sprintf("%x", id[:])
}

func unpack(r *http.Request) *RequestVars {
	vars := mux.Vars(r)
	rv := &RequestVars{
		Oid:           vars["oid"],
		Authorization: r.Header.Get("Authorization"),
	}

	if r.Method == "POST" { // Maybe also check if +json
		var p RequestVars
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			return rv
		}

		rv.Oid = p.Oid
		rv.Size = p.Size
	}

	return rv
}

// TODO cheap hack, unify with unpack
func unpackBatch(r *http.Request) *BatchVars {
	var bv BatchVars

	err := json.NewDecoder(r.Body).Decode(&bv)
	if err != nil {
		return &bv
	}

	for i := 0; i < len(bv.Objects); i++ {
		bv.Objects[i].Authorization = r.Header.Get("Authorization")
	}

	return &bv
}

func writeStatus(w http.ResponseWriter, r *http.Request, status int) {
	message := http.StatusText(status)

	mediaParts := strings.Split(r.Header.Get("Accept"), ";")
	mt := mediaParts[0]
	if strings.HasSuffix(mt, "+json") {
		message = `{"message":"` + message + `"}`
	}

	w.WriteHeader(status)
	fmt.Fprint(w, message)
}
