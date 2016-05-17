package feed_line_repo

import (
	"time"

	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gopkg.in/gorp.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//-----------------------------------------------------------------------------------//

//Feedline Input Model
type feedLineInputModel struct {
	models.FeedLineUnit `bson:",inline"`
	RetryCount          uint                `bson:"retry_count"`
	Status              feedLineQueueStatus `bson:"status"`
	IdString            string              `bson:"id_string"`
	MacroTaskIdString   string              `bson:"macro_task_id_string"`
}

//-----------------------------------------------------------------------------------//

//Feedline queue status
type feedLineQueueStatus uint

const (
	queued feedLineQueueStatus = iota
	retry
	success
	failed
)

//-----------------------------------------------------------------------------------//

//Client to get/set FeedLineUnits from mongodb
type inputQueue struct {
	mgo *mgo.Database
}

func (i *inputQueue) Add(flu models.FeedLineUnit) (id uuid.UUID, err error) {

	switch {
	case flu.ID == uuid.Nil:
		flu.ID = uuid.NewV4()
		fallthrough
	case flu.CreatedAt.Valid == false:
		flu.CreatedAt = gorp.NullTime{time.Now(), true}
	}

	err = i.mgo.C("feedline_input").
		Insert(&feedLineInputModel{
			FeedLineUnit:      flu,
			RetryCount:        0,
			Status:            queued,
			IdString:          flu.ID.String(),
			MacroTaskIdString: flu.MacroTaskId.String(),
		})
	if err != nil {
		//Error Code from mgo
		if err.(*mgo.LastError).Code == 11000 {
			return uuid.Nil, ErrDuplicateReferenceId
		}
	}
	return flu.ID, err
}

func (i *inputQueue) Get(id uuid.UUID) (flu models.FeedLineUnit, err error) {

	err = i.mgo.C("feedline_input").
		FindId(id).
		One(&flu)
	if err != nil && err == mgo.ErrNotFound {
		err = ErrFLUNotFoundInInputQueue
	}
	return flu, err
}

func (i *inputQueue) GetQueuedRaw() ([]feedLineInputModel, error) {

	var flus []feedLineInputModel
	err := i.mgo.C("feedline_input").
		Find(bson.M{"status": queued}).
		All(&flus)
	return flus, err
}

func (i *inputQueue) GetQueued() ([]models.FeedLineUnit, error) {

	var flus []models.FeedLineUnit
	err := i.mgo.C("feedline_input").
		Find(bson.M{"status": queued}).
		All(&flus)

	return flus, err
}

func (i *inputQueue) MarkFinished() error {

	_, err := i.mgo.C("feedline_input").UpdateAll(bson.M{"status": queued}, bson.M{"$set": bson.M{"status": success}})
	return err
}
