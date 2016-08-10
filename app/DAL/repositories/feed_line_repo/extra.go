package feed_line_repo

import (
	"github.com/crowdflux/angel/app/DAL/clients/postgres"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"gopkg.in/mgo.v2/bson"
)

func SyncAll() {

	inpQ := NewInputQueue()

	existingQFlus := []feedLineInputModel{}

	err := inpQ.mgo.C("feedline_input").Find(bson.M{}).All(&existingQFlus)
	if err != nil {
		plog.Error("feedline", err)
		return
	}

	plog.Info("existingQFlus", len(existingQFlus))

	fluRepo := fluRepo{
		Db: postgres.GetPostgresClient(),
	}

	flus := []models.FeedLineUnit{}

	_, err = fluRepo.Db.Select(&flus, "SELECT * FROM feed_line")
	if err != nil {
		plog.Error("feedline", err)
		return
	}

	plog.Info("flus", len(flus))

	//flusToUpdate := []feedLineInputModel{}
	flusToInsert := []interface{}{}

	for _, eflu := range existingQFlus {

		if _, ok := existsInList(flus, eflu.ID); !ok {

			f := eflu.FeedLineUnit

			f.Build = models.JsonF{
				"fail": true,
			}
			flusToInsert = append(flusToInsert, &f)
		}
	}

	//for _, f := range flusToInsert {
	//	plog.Info("flusToInsert", f.(*models.FeedLineUnit).ID)
	//
	//}

	err = fluRepo.Db.Insert(flusToInsert...)
	if err != nil {
		plog.Error("feedline", err)
		return
	}

	plog.Info("extra", len(flusToInsert), " new flus inserted")
}

func SyncAllBack() {

	inpQ := NewInputQueue()

	existingQFlus := []feedLineInputModel{}

	err := inpQ.mgo.C("feedline_input").Find(bson.M{}).All(&existingQFlus)
	if err != nil {
		plog.Error("feedline", err)
		return
	}

	plog.Info("existingQFlus", len(existingQFlus))

	fluRepo := fluRepo{
		Db: postgres.GetPostgresClient(),
	}

	flus := []models.FeedLineUnit{}

	_, err = fluRepo.Db.Select(&flus, "SELECT * FROM feed_line")
	if err != nil {
		plog.Error("feedline", err)
		return
	}

	plog.Info("flus", len(flus))

	flusToUpdate := []feedLineInputModel{}
	flusToInsert := []interface{}{}

	for _, eflu := range flus {

		if f, ok := existsInListInQ(existingQFlus, eflu.ID); ok {

			if f.Status == queued {

				flusToUpdate = append(flusToUpdate, *f)
			}
		} else {

			//flusToInsert = append(flusToInsert, )
		}
	}

	err = fluRepo.Db.Insert(flusToInsert...)
	if err != nil {
		plog.Error("feedline", err)
		return
	}

	plog.Info("extra", len(flusToInsert), " new flus inserted")

	if len(flusToInsert) > 0 {
		err = inpQ.mgo.C("feedline_input").Insert(flusToInsert...)
		if err != nil {
			plog.Error("feedline", err)
		}

	}

	for _, updateFlu := range flusToUpdate {

		err = inpQ.mgo.C("feedline_input").UpdateId(updateFlu.ID, updateFlu)
		if err != nil {
			plog.Error("feedline", err, updateFlu)
		}
	}
}

func existsInList(list []models.FeedLineUnit, toFindId uuid.UUID) (*models.FeedLineUnit, bool) {
	for _, elem := range list {
		if elem.ID == toFindId {

			plog.Info("asd", elem.ID)
			return &elem, true
		}
	}
	return nil, false
}

func existsInListInQ(list []feedLineInputModel, toFindId uuid.UUID) (*feedLineInputModel, bool) {
	for _, elem := range list {
		if elem.ID == toFindId {

			plog.Info("asd", elem.ID)
			return &elem, true
		}
	}
	return nil, false
}
