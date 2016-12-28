package flu_monitor

import (
	"database/sql"
	"encoding/json"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/services/flu_logger_svc"
)

func putDbLog(flusToLog []feed_line.FLU, message string, resp WebhookResponse) {

	dbLogArr := make([]models.FeedLineLog, len(flusToLog))
	jsObj := models.JsonF{}
	jsonBytes, _ := json.Marshal(resp)
	jsObj.Scan(string(jsonBytes))
	for i, fl := range flusToLog {
		dbLog := models.FeedLineLog{
			//ID         int            `db:"id" json:"id" bson:"_id"`
			FluId:       fl.ID,
			Message:     sql.NullString{message, true},
			MetaData:    jsObj,
			Event:       10,
			StepType:    sql.NullInt64{int64(12), true},
			StepId:      fl.StepId,
			CreatedAt:   fl.CreatedAt,
			MasterFluId: fl.MasterId,
		}
		dbLogArr[i] = dbLog
	}

	flu_logger_svc.LogRaw(dbLogArr)
}
