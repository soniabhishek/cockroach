package clients

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

func TestCrowdy(t *testing.T) {

	t.SkipNow()

	cc := crowdyClient{}

	flu := models.FeedLineUnit{
		ID: uuid.NewV4(),
	}

	success, err := cc.PushFLU(flu)

	assert.True(t, success)
	assert.NoError(t, err)
}
