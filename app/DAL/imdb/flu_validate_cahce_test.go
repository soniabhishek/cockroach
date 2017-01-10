package imdb

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/crowdflux/angel/app/models"

	"github.com/crowdflux/angel/app/models/uuid"
)

func TestEvalExpressCache_Get(t *testing.T) {
	idb := NewFluValidateCache()
	val := []models.FLUValidator{
		models.FLUValidator{
			FieldName:   "brand",
			Type:        "STRING",
			IsMandatory: true,
		},
		models.FLUValidator{
			FieldName:   "color",
			Type:        "STRING",
			IsMandatory: true,
		},
		models.FLUValidator{
			FieldName:   "image_url",
			Type:        "IMAGE_ARRAY",
			IsMandatory: false,
		},
		models.FLUValidator{
			FieldName:   "image_single",
			Type:        "IMAGE",
			IsMandatory: false,
		},
		models.FLUValidator{
			FieldName:   "category_id",
			Type:        "STRING",
			IsMandatory: true,
		}}
	key := uuid.NewV4()
	idb.Set(key, val)
	res, err := idb.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, val, res)
}
