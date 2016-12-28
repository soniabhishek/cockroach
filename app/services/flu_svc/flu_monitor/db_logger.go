package flu_monitor

import (
	"github.com/crowdflux/angel/app/models"
	"encoding/json"
	"database/sql"
	"github.com/crowdflux/angel/app/services/flu_logger_svc"
)

func putDbLog(completedFLUs []models.FeedLineUnit, message string, resp FluResponse) {

	dbLogArr := make([]models.FeedLineLog, len(completedFLUs))
	jsObj := models.JsonF{}
	jsonBytes, _ := json.Marshal(resp)
	jsObj.Scan(string(jsonBytes))
	for i, fl := range completedFLUs {
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