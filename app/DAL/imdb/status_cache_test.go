package imdb

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/docker/libnetwork/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFluUploadCache(t *testing.T) {
	idb := newFluUploadCache()
	data := models.FluUploadStats{}
	data.Status = 2
	idb.Set("hello", data)
	val, err := idb.Get("hello")
	assert.NoError(t, err)
	assert.Equal(t, val, data)
	val, err = idb.SafeSet("hello", data)
	assert.Error(t, err)
}
