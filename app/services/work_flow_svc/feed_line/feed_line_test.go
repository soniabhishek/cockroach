package feed_line_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/feed_line"
)

func TestNew(t *testing.T) {

	fl := feed_line.New()

	fluId := uuid.NewV4()

	fl <- feed_line.FLU{
		FeedLineUnit: models.FeedLineUnit{
			ID: fluId,
		},
	}

	flu := <-fl

	assert.EqualValues(t, fluId, flu.ID)

}
