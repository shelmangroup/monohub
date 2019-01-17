package lfs

import (
	"time"
)

// RequestVars contain variables from the HTTP request. Variables from routing, json body decoding, and
// some headers are stored.
type RequestVars struct {
	Oid           string
	Size          int64
	Repo          string
	Authorization string
}

type BatchVars struct {
	Transfers []string       `json:"transfers,omitempty"`
	Operation string         `json:"operation"`
	Objects   []*RequestVars `json:"objects"`
}

// MetaObject is object metadata as seen by the object and metadata stores.
type MetaObject struct {
	Oid      string `json:"oid"`
	Size     int64  `json:"size"`
	Existing bool
}

type BatchResponse struct {
	Transfer string            `json:"transfer,omitempty"`
	Objects  []*Representation `json:"objects"`
}

// Representation is object medata as seen by clients of the lfs server.
type Representation struct {
	Oid     string           `json:"oid"`
	Size    int64            `json:"size"`
	Actions map[string]*link `json:"actions"`
	Error   *ObjectError     `json:"error,omitempty"`
}

type ObjectError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type User struct {
	Name string `json:"name"`
}

type Lock struct {
	Id       string    `json:"id"`
	Path     string    `json:"path"`
	Owner    User      `json:"owner"`
	LockedAt time.Time `json:"locked_at"`
}

type LockRequest struct {
	Path string `json:"path"`
}

type LockResponse struct {
	Lock    *Lock  `json:"lock"`
	Message string `json:"message,omitempty"`
}

type UnlockRequest struct {
	Force bool `json:"force"`
}

type UnlockResponse struct {
	Lock    *Lock  `json:"lock"`
	Message string `json:"message,omitempty"`
}

type LockList struct {
	Locks      []Lock `json:"locks"`
	NextCursor string `json:"next_cursor,omitempty"`
	Message    string `json:"message,omitempty"`
}

type VerifiableLockRequest struct {
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

type VerifiableLockList struct {
	Ours       []Lock `json:"ours"`
	Theirs     []Lock `json:"theirs"`
	NextCursor string `json:"next_cursor,omitempty"`
	Message    string `json:"message,omitempty"`
}

// link provides a structure used to build a hypermedia representation of an HTTP link.
type link struct {
	Href      string            `json:"href"`
	Header    map[string]string `json:"header,omitempty"`
	ExpiresAt time.Time         `json:"expires_at,omitempty"`
}
