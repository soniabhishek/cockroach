package feed_line

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
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
