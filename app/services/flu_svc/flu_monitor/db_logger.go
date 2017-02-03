package flu_monitor

import (
	"database/sql"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/services/flu_logger_svc"
	"net/http"
)

func putDbLog(flusToLog map[uuid.UUID]feed_line.FLU, resp WebhookResponse) {

	dbLogArr := make([]models.FeedLineLog, len(flusToLog))
	jsObj := models.JsonF{}
	jsonBytes, _ := json.Marshal(resp)
	jsObj.Scan(string(jsonBytes))

	invalidFlusMap := make(map[uuid.UUID]invalidFlu)
	for _, flu := range resp.Invalid_Flus {
		flId := uuid.FromStringOrNil(flu.FluID)
		if flId != uuid.Nil {
			invalidFlusMap[flId] = flu
		}
	}

	var message string
	i := 0
	for _, fl := range flusToLog {
		metadata := jsObj.Copy()

		if resp.HttpStatusCode == http.StatusOK {
			message = "SUCCESS"
		} else {
			message = "ERROR"
		}
		val, ok := invalidFlusMap[fl.ID]
		if ok {
			message = "INVALID FLU"
			metadata.Merge(models.JsonF{"Error": val.Error, "Message": val.Message})
		}

		dbLog := models.FeedLineLog{
			//ID         int            `db:"id" json:"id" bson:"_id"`
			FluId:       fl.ID,
			Message:     sql.NullString{message, true},
			MetaData:    metadata,
			Event:       10,
			StepType:    sql.NullInt64{int64(12), true},
			StepId:      fl.StepId,
			CreatedAt:   fl.CreatedAt,
			MasterFluId: fl.MasterId,
		}
		dbLogArr[i] = dbLog
		i++
	}

	flu_logger_svc.LogRaw(dbLogArr)
}

func putDbLogCustom(flusToLog map[uuid.UUID]feed_line.FLU, message string, metadata models.JsonF) {

	dbLogArr := make([]models.FeedLineLog, len(flusToLog))
	i := 0
	for _, fl := range flusToLog {

		dbLog := models.FeedLineLog{
			//ID         int            `db:"id" json:"id" bson:"_id"`
			FluId:       fl.ID,
			Message:     sql.NullString{message, true},
			MetaData:    metadata,
			Event:       10,
			StepType:    sql.NullInt64{int64(12), true},
			StepId:      fl.StepId,
			CreatedAt:   fl.CreatedAt,
			MasterFluId: fl.MasterId,
		}
		dbLogArr[i] = dbLog
		i++
	}

	flu_logger_svc.LogRaw(dbLogArr)
}
