package feed_line

import (
	"testing"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	fl := make(Fl, 100)

	fluId := uuid.NewV4()

	fl <- FLU{
		FeedLineUnit: models.FeedLineUnit{
			ID: fluId,
		},
	}

	flu := <-fl

	assert.EqualValues(t, fluId, flu.ID)

}
