package imdb

import (
	"github.com/docker/libnetwork/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"testing"
)

func TestInit(t *testing.T) {
	idb := new()
	idb.Set("hello", "abhishek")
	val, err := idb.Get("hello")
	assert.NoError(t, err)
	assert.Equal(t, val, "abhishek")
	val, err = idb.SafeSet("hello", "stranger")
	assert.Error(t, err)
}
