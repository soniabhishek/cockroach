package crowdsourcing_step

import (
	"errors"
	"gitlab.com/playment-main/angel/app/DAL/repositories/feed_line_repo"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/plog"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/feed_line"
)

type FluUpdate struct {
	FluId       uuid.UUID       `json:"flu_id"`
	BuildUpdate models.JsonFake `json:"build_update"`
}

func FluUpdateHandler(updates []FluUpdate) error {

	flus := Std.GetBuffered()
	flr := feed_line_repo.New()

	updatable := []feed_line.FLU{}

	for _, update := range updates {

		flu, ok := flus[update.FluId]

		if !ok {
			// Handle error
			plog.Error("Flu Handler", errors.New("Flu Not present in the buffer"), update.FluId)
			continue
		}

		flu.Build.Merge(update.BuildUpdate)

		updatable = append(updatable, flu)

	}

	feedLineUnits := []models.FeedLineUnit{}

	for _, flu := range updatable {
		feedLineUnits = append(feedLineUnits, flu.FeedLineUnit)
	}

	err := flr.BulkUpdate(feedLineUnits)
	if err != nil {
		plog.Error("Flu Handler Bulk Update, Aborting", err)
		return err
	}

	for _, flu := range updatable {
		ok := Std.finishFlu(flu)
		if !ok {
			plog.Error("Flu Handler", errors.New("finishFlu false for "+flu.ID.String()))
		}
	}

	return nil
}
