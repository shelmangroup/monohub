package storage

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepoPath(t *testing.T) {
	r := NewGitRepository("")
	assert.Equal(t, r.Path, ".")

	r = NewGitRepository("/foo/../bar")
	assert.Equal(t, r.Path, "/bar")
}

func TestInitialize(t *testing.T) {
	dir, err := ioutil.TempDir("", "monohub-test-initialize")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	r := NewGitRepository(dir)
	assert.Equal(t, r.Path, dir)

	err = r.Init()
	assert.Nil(t, err, "A new git repository should be created on init if none exists")

	err = r.Init()
	assert.Nil(t, err, "Initializing an existing repo should be a noop")
}
