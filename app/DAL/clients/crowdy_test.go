package clients

//
//import (
//	"github.com/stretchr/testify/assert"
//	"gitlab.com/playment-main/angel/app/DAL/clients/postgres"
//	"gitlab.com/playment-main/angel/app/DAL/repositories/feed_line_repo"
//	"gitlab.com/playment-main/angel/app/models"
//	"gitlab.com/playment-main/angel/app/models/uuid"
//	"testing"
//)
//
//func TestCrowdy(t *testing.T) {
//
//	cc := crowdyClient{}
//
//	var step models.Step
//
//	postgres.GetPostgresClient().SelectOne()
//
//	flu := models.FeedLineUnit{
//		ID:     uuid.FromStringOrNil("25cdce1c-f3f4-4ee3-bc3f-b0ad2afc52c5"),
//		StepId: uuid.FromStringOrNil("2b26a671-d635-489d-8da5-7df3c0d29f2a"),
//	}
//
//	feed_line_repo.New().Save()
//
//	success, err := cc.PushFLU(flu)
//
//	assert.True(t, success)
//	assert.NoError(t, err)
//}
