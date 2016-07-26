package crowdsourcing_step

import (
	"errors"

	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/work_flow_svc/feed_line"
)

type FluUpdate struct {
	FluId       uuid.UUID    `json:"flu_id"`
	BuildUpdate models.JsonF `json:"build_update"`
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
