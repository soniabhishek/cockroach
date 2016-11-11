package clients

import (
	"testing"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCrowdy(t *testing.T) {

	t.SkipNow()

	cc := crowdyClient{}

	flu := models.FeedLineUnit{
		ID: uuid.NewV4(),
	}

	micId := uuid.NewV4()

	success, err := cc.PushFLU(flu, micId)

	assert.True(t, success)
	assert.NoError(t, err)
}
