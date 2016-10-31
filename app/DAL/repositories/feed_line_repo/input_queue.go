package feed_line_repo

import (
	"time"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/lib/pq"
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
	ProjectIdString     string              `bson:"project_id_string"`
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
		flu.CreatedAt = pq.NullTime{time.Now(), true}
	}

	flu.IsMaster = true
	flu.IsActive = true
	flu.MasterId = flu.ID

	err = i.mgo.C("feedline_input").
		Insert(&feedLineInputModel{
			FeedLineUnit:    flu,
			RetryCount:      0,
			Status:          queued,
			IdString:        flu.ID.String(),
			ProjectIdString: flu.ProjectId.String(),
		})
	if err != nil {
		// Error Code from mgo for duplicate id
		// unsafe
		if e, ok := err.(*mgo.LastError); ok && e.Code == 11000 {
			return uuid.Nil, ErrDuplicateReferenceId
		} else {
			plog.Error("Input queue", err, "mongo insert failed")
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

func (i *inputQueue) MarkFinished(flus []models.FeedLineUnit) error {

	fluIdsString := make([]string, len(flus))
	for i, flu := range flus {
		fluIdsString[i] = flu.ID.String()
	}

	_, err := i.mgo.C("feedline_input").UpdateAll(bson.M{"id_string": bson.M{"$in": fluIdsString}}, bson.M{"$set": bson.M{"status": success}})
	return err
}
