package flu_logger_svc

import (
	"database/sql"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/lib/pq"
	"time"
)

func LogStepEntry(flu models.FeedLineUnit, stepType step_type.StepType) {

	log(flu, true, stepType)
}

func LogStepExit(flu models.FeedLineUnit, stepType step_type.StepType) {

	log(flu, false, stepType)
}

func log(flu models.FeedLineUnit, stepEntry bool, stepType step_type.StepType) {

	fluLog := models.FeedLineLog{
		FluId:      flu.ID,
		Message:    sql.NullString{"", false},
		MetaData:   flu.Build,
		StepType:   sql.NullInt64{int64(stepType), true},
		StepEntry:  sql.NullBool{stepEntry, true},
		StepExit:   sql.NullBool{!stepEntry, true},
		StepId:     flu.StepId,
		WorkFlowId: uuid.Nil,
		CreatedAt:  pq.NullTime{time.Now(), true},
	}

	feed_line.GetFeedlineLoggerChannel().Push(fluLog)
}
