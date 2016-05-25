package flu_output

import (
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"fmt"
)

var feedLinePipe = make(map[uuid.UUID][]models.FeedLineUnit)

type FluMonitor struct {

}

func (fm *FluMonitor) AddToOutputQueue(flu models.FeedLineUnit) error {
	feedLineArr := make([]models.FeedLineUnit, 1)
	feedLineArr[0] = flu
	return fm.AddManyToOutputQueue(feedLineArr)
}

func (fm *FluMonitor) AddManyToOutputQueue(fluBundle []models.FeedLineUnit) error {

	for _, flu:=range fluBundle{
		feedLinePipe[flu.ID] = append(feedLinePipe[flu.ID], flu)
	}

	var sendersIds = make([]uuid.UUID, 1)
	for client_id:=range feedLinePipe{
		stores := feedLinePipe[client_id]
		length := len(stores)
		fmt.Println(client_id, feedLinePipe[client_id])
		if length >= 1000{ //some threshold value TODO read from config
			sendersIds = append(sendersIds, client_id)
		}
	}
	return
}
