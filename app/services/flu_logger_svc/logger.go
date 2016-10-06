package flu_logger_svc

import (
	"database/sql"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/lib/pq"
	"time"
)

func LogStepEntry(flu models.FeedLineUnit, stepType step_type.StepType, retried bool) {

	log(flu, 1, stepType, "", retried)
}

func LogStepExit(flu models.FeedLineUnit, stepType step_type.StepType, retried bool) {

	log(flu, 2, stepType, "", retried)
}

func LogStepError(flu models.FeedLineUnit, stepType step_type.StepType, errMsg string, retried bool) {
	log(flu, 3, stepType, errMsg, retried)
}

func log(flu models.FeedLineUnit, event int, stepType step_type.StepType, message string, retried bool) {

	metaData := models.JsonF{"build": flu.Build}
	if retried {
		metaData["retry"] = true
	}

	fluLog := models.FeedLineLog{
		FluId:     flu.ID,
		Message:   sql.NullString{message, !IsEmpty(message)},
		MetaData:  metaData,
		Event:     event,
		StepType:  sql.NullInt64{int64(stepType), true},
		StepId:    flu.StepId,
		CreatedAt: pq.NullTime{time.Now(), true},
	}

	feed_line.GetFeedlineLoggerChannel().Push(fluLog)
}

func IsEmpty(s string) bool {
	if s == "" {
		return true
	}
	return false
}
