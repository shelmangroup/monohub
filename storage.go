package main

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/iterator"
	"io"
	"io/ioutil"

	gcs "cloud.google.com/go/storage"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/format/index"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
	"gopkg.in/src-d/go-git.v4/storage"
)

var (
	gcsBucket = kingpin.Flag("gcs-bucket", "Google Cloud Storage Bucket.").Required().String()
)

type Storage struct {
	bucket *gcs.BucketHandle
}

type EncodedObjectIter struct {
	storage *Storage
	t       plumbing.ObjectType
	objs    *gcs.ObjectIterator
}

func NewStorage() (*Storage, error) {
	log.WithFields(log.Fields{
		"bucket": *gcsBucket,
	}).Debug("Creating GCS session")

	client, err := gcs.NewClient(context.Background())
	if err != nil {
		log.Errorf("Failed to create gcs client: %s", err.Error())
		return nil, err
	}
	bh := client.Bucket(*gcsBucket)

	return &Storage{
		bucket: bh,
	}, nil
}

func (s *Storage) NewEncodedObject() plumbing.EncodedObject {
	log.Debug("NewEncodedObject")
	return &plumbing.MemoryObject{}
}

func (s *Storage) SetEncodedObject(obj plumbing.EncodedObject) (plumbing.Hash, error) {
	log.Debug("SetEncodedObject")
	ctx := context.Background()

	key := s.buildObjectKey(obj.Hash(), obj.Type())

	r, err := obj.Reader()
	if err != nil {
		return obj.Hash(), err
	}

	wc := s.bucket.Object(key).NewWriter(ctx)
	if _, err = io.Copy(wc, r); err != nil {
		return obj.Hash(), err
	}
	if err := wc.Close(); err != nil {
		return obj.Hash(), err
	}

	return obj.Hash(), nil
}

func (s *Storage) buildObjectKey(h plumbing.Hash, t plumbing.ObjectType) string {
	hash := h.String()
	return fmt.Sprintf("objects/%s/%s/%s", hash[0:2], hash[2:4], hash)
}

func (s *Storage) EncodedObject(t plumbing.ObjectType, h plumbing.Hash) (plumbing.EncodedObject, error) {
	log.Debug("EncodedObject")
	var err error
	ctx := context.Background()

	key := s.buildObjectKey(h, t)

	r, err := s.bucket.Object(key).NewReader(ctx)
	if err == gcs.ErrObjectNotExist {
		return nil, plumbing.ErrObjectNotFound
	}
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return objectFromRecord(data, t)
}

func (s *Storage) HasEncodedObject(h plumbing.Hash) error {
	log.Debug("HasEncodedObject")
	ctx := context.Background()

	key := s.buildObjectKey(h, plumbing.AnyObject)
	_, err := s.bucket.Object(key).Attrs(ctx)
	return err
}

func (s *Storage) EncodedObjectSize(h plumbing.Hash) (int64, error) {
	log.Debug("EncodedObjectSize")
	ctx := context.Background()

	key := s.buildObjectKey(h, plumbing.AnyObject)

	a, err := s.bucket.Object(key).Attrs(ctx)
	if err != nil {
		return 0, err
	}

	return a.Size, nil
}

func objectFromRecord(data []byte, t plumbing.ObjectType) (plumbing.EncodedObject, error) {
	o := &plumbing.MemoryObject{}
	o.SetType(t)
	o.SetSize(int64(len(data)))

	_, err := o.Write(data)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (s *Storage) IterEncodedObjects(t plumbing.ObjectType) (storer.EncodedObjectIter, error) {
	ctx := context.Background()
	objs := s.bucket.Objects(ctx, &gcs.Query{Prefix: fmt.Sprintf("objects/%s/", t.String())})
	return &EncodedObjectIter{
		storage: s,
		t:       t,
		objs:    objs,
	}, nil
}

func (i *EncodedObjectIter) Next() (plumbing.EncodedObject, error) {
	ctx := context.Background()

	a, err := i.objs.Next()
	if err != nil {
		return nil, err
	}

	r, err := i.storage.bucket.Object(a.Name).NewReader(ctx)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return objectFromRecord(data, i.t)
}

func (i *EncodedObjectIter) ForEach(cb func(obj plumbing.EncodedObject) error) error {
	for {
		obj, err := i.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		if err := cb(obj); err != nil {
			if err == storer.ErrStop {
				return nil
			}

			return err
		}
	}
}

func (i *EncodedObjectIter) Close() {}

func (s *Storage) Reference(n plumbing.ReferenceName) (*plumbing.Reference, error) {
	log.Debug("Reference")
	ctx := context.Background()

	key := s.buildReferenceKey(n)

	r, err := s.bucket.Object(key).NewReader(ctx)
	if err == gcs.ErrObjectNotExist {
		return nil, plumbing.ErrReferenceNotFound
	}
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return plumbing.NewReferenceFromStrings(
		n.String(),
		string(data),
	), nil
}

func (s *Storage) SetReference(ref *plumbing.Reference) error {
	log.Debug("SetReference")
	ctx := context.Background()

	key := s.buildReferenceKey(ref.Name())

	w := s.bucket.Object(key).NewWriter(ctx)

	_, err := fmt.Fprint(w, ref.Target())
	if err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) RemoveReference(n plumbing.ReferenceName) error {
	log.Debug("RemoveReference")
	ctx := context.Background()
	key := s.buildReferenceKey(n)
	return s.bucket.Object(key).Delete(ctx)
}

func (s *Storage) CheckAndSetReference(ref, old *plumbing.Reference) error {
	log.Debug("CheckAndSetReference")
	return s.SetReference(ref) // FIXME: nope!
}

func (s *Storage) buildReferenceKey(n plumbing.ReferenceName) string {
	log.Debug("buildReferenceKey")
	return n.String()
}

func (s *Storage) IterReferences() (storer.ReferenceIter, error) {
	log.Debug("IterReferences")
	ctx := context.Background()
	objs := s.bucket.Objects(ctx, &gcs.Query{Prefix: "refs/"})

	var refs []*plumbing.Reference
	for {
		o, err := objs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		r, err := s.bucket.Object(o.Name).NewReader(ctx)
		if err != nil {
			return nil, err
		}

		target, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}

		refs = append(refs, plumbing.NewReferenceFromStrings(o.Name, string(target)))
	}

	return storer.NewReferenceSliceIter(refs), nil
}

func (s *Storage) CountLooseRefs() (int, error) {
	log.Debug("CountLooseRefs")
	return 0, nil // FIXME: wtf?
}

func (s *Storage) PackRefs() error {
	log.Debug("PackRefs")
	return nil
}

func (s *Storage) Config() (*config.Config, error) {
	log.Debug("Config")
	ctx := context.Background()

	key := s.buildConfigKey()
	r, err := s.bucket.Object(key).NewReader(ctx)
	if err == gcs.ErrObjectNotExist {
		return config.NewConfig(), nil
	}
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	c := &config.Config{}
	return c, json.Unmarshal(data, c)
}

func (s *Storage) SetConfig(r *config.Config) error {
	log.Debug("SetConfig")
	ctx := context.Background()

	key := s.buildConfigKey()

	json, err := json.Marshal(r)
	if err != nil {
		return err
	}

	w := s.bucket.Object(key).NewWriter(ctx)
	_, err = w.Write(json)
	if err != nil {
		return err
	}

	return w.Close()
}

func (s *Storage) buildConfigKey() string {
	return "config"
}

func (s *Storage) Index() (*index.Index, error) {
	log.Debug("Index")
	ctx := context.Background()

	key := s.buildIndexKey()

	r, err := s.bucket.Object(key).NewReader(ctx)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	idx := &index.Index{}
	return idx, json.Unmarshal(data, idx)
}

func (s *Storage) SetIndex(idx *index.Index) error {
	log.Debug("SetIndex")
	ctx := context.Background()

	key := s.buildIndexKey()

	json, err := json.Marshal(idx)
	if err != nil {
		return err
	}

	w := s.bucket.Object(key).NewWriter(ctx)
	_, err = w.Write(json)
	if err != nil {
		return err
	}

	return w.Close()
}

func (s *Storage) buildIndexKey() string {
	return "index"
}

func (s *Storage) Shallow() ([]plumbing.Hash, error) {
	log.Debug("Shallow")
	ctx := context.Background()

	key := s.buildShallowKey()

	r, err := s.bucket.Object(key).NewReader(ctx)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var h []plumbing.Hash
	return h, json.Unmarshal(data, h)
}

func (s *Storage) SetShallow(hash []plumbing.Hash) error {
	log.Debug("SetShallow")
	ctx := context.Background()

	key := s.buildShallowKey()

	json, err := json.Marshal(hash)
	if err != nil {
		return err
	}

	w := s.bucket.Object(key).NewWriter(ctx)
	_, err = w.Write(json)
	if err != nil {
		return err
	}

	return w.Close()
}

func (s *Storage) buildShallowKey() string {
	return "shallow"
}

func (s *Storage) Module(name string) (storage.Storer, error) {
	log.Debug("Module")
	return nil, nil
}
