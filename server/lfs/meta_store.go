package lfs

import (
	"github.com/boltdb/bolt"
)

// MetaStore implements a metadata storage. It stores user credentials and Meta
// information for objects. The storage is handled by boltdb.
type MetaStore struct {
	db *bolt.DB
}
