package crowdsourcing_step_svc

import (
	"errors"

	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/plog/log_tags"
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
			plog.Error("Flu Handler crowdy", errors.New("Flu Not present in the buffer"), plog.MessageWithParam(log_tags.FLU_ID, update.FluId))
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
		plog.Error("Flu Handler Crowdy", err, plog.Message("Flu Handler Bulk Update, Aborting"))
		return err
	}

	for _, flu := range updatable {
		ok := Std.finishFlu(flu)
		if !ok {
			plog.Error("Flu Handler Crowdy", errors.New("finishFlu false for "+flu.ID.String()))
		}
	}

	return nil
}

func FluUpdateHandlerCustom(updates []FluUpdate) error {

	flr := feed_line_repo.New()

	flus := []models.FeedLineUnit{}

	for _, update := range updates {
		flus = append(flus, models.FeedLineUnit{
			ID:    update.FluId,
			Build: update.BuildUpdate,
		})
	}

	updatedFlus, nonUpdableFlus, err := flr.BulkFluBuildUpdateByStepType(flus, step_type.CrowdSourcing)
	if err != nil {
		if err != feed_line_repo.ErrPartiallyUpdatedFlus && err != feed_line_repo.ErrNoUpdatableFlus {
			plog.Error("Flu Handler Crowdy", err, plog.Message("Flu Handler Bulk Update, Aborting"))
			return err
		} else {

			nonUpdatableIds := []uuid.UUID{}
			for _, flu := range nonUpdableFlus {
				nonUpdatableIds = append(nonUpdatableIds, flu.ID)
			}

			plog.Error("Flu Handler Crowdy", err, plog.Message("Crowdy flu handler partially updated. Non_Updated flu_ids"), plog.MessageWithParam(log_tags.FLU_ID, nonUpdatableIds))
			// this wont return
			// this will continue
		}
	}

	go func() {

		for _, flu := range updatedFlus {
			ok := Std.finishFlu(feed_line.FLU{FeedLineUnit: flu})
			if !ok {
				plog.Error("Flu Handler Crowdy", errors.New("finishFlu false for "+flu.ID.String()))
			}
		}
	}()

	// it will return nil in case of no error
	// return ErrPartiallyUpdatedFlus in case of partial update
	return err
}
