package uuid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromCEnc(t *testing.T) {

	id := NewV4()

	s := toCEnc(id)

	nid, err := FromCEnc(s)
	assert.NoError(t, err)
	assert.EqualValues(t, id.String(), nid.String())
	assert.EqualValues(t, id, nid)
	assert.NotContains(t, s, "=")

}
