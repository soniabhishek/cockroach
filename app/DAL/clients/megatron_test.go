package clients

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMegatron(t *testing.T) {

	t.Skip("Create template before transformation")

	cc := megatronClient{}
	testJson := models.JsonF{"id": 1, "template": "template"}
	transformedData, err := cc.Transform(testJson, "57a851f7a1535f9f6cd716af")
	assert.NoError(t, err)
	for key, val := range testJson {
		assert.EqualValues(t, val, transformedData[key])
	}
}
