package crowdsourcing_step

import (
	"gitlab.com/playment-main/angel/app/DAL/repositories/feed_line_repo"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

type FluUpdate struct {
	FluId       uuid.UUID       `json:"flu_id"`
	BuildUpdate models.JsonFake `json:"build_update"`
}

func FluUpdateHandler(updates []FluUpdate) error {

	flus := Std.GetBuffered()

	for _, update := range updates {

		flu := flus.Get(update.FluId)

		flu.Build = MergeJsonFake(flu.Build, update.BuildUpdate)

		feed_line_repo.New().Save(flu.FeedLineUnit)

		Std.finishFlu(flu)
	}
	return nil
}

func MergeJsonFake(a models.JsonFake, b models.JsonFake) (c models.JsonFake) {
	for k, v := range a {
		b[k] = v
	}
	return b
}
