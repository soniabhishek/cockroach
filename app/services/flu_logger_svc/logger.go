package feed_line

import (
	"database/sql"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/services/work_flow_svc/feed_line"
	"github.com/lib/pq"
	"time"
)

func LogStepEntry(flu feed_line.FLU, metaData ...models.JsonF) {

	log(flu, true, metaData)
}

func LogStepExit(flu feed_line.FLU, metaData ...models.JsonF) {

	log(flu, false, metaData)
}

func log(flu feed_line.FLU, stepEntry bool, metaData ...models.JsonF) {

	metaDataMerged := models.JsonF{}

	for _, m := range metaData {
		metaDataMerged.Merge(m)
	}

	metaDataMerged.Merge(models.JsonF{"build": flu.Build})

	fluLog := models.FeedLineLog{
		FluId:      flu.ID,
		Message:    sql.NullString{"asfas", true},
		MetaData:   metaDataMerged,
		StepType:   sql.NullInt64{-1, false},
		StepEntry:  sql.NullBool{stepEntry, true},
		StepExit:   sql.NullBool{!stepEntry, true},
		StepId:     flu.StepId,
		WorkFlowId: uuid.Nil,
		CreatedAt:  pq.NullTime{time.Now(), true},
	}
}
