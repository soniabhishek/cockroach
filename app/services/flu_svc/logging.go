package flu_svc

import (
	"log"
	"time"

	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

// This is an example implementation of logger middleware
// DO NOT USE
func loggingMiddleware(logger log.Logger) func(IFluService) IFluService {
	return func(next IFluService) IFluService {
		return logmw{logger, next}
	}
}

type logmw struct {
	logger log.Logger
	IFluService
}

func (mw logmw) AddFeedLineUnit(flu *models.FeedLineUnit) (err error) {
	defer func(begin time.Time) {

		mw.logger.Print(
			"method", "AddFeedLineUnit",
			"input", flu,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.IFluService.AddFeedLineUnit(flu)
	return
}

func (mw logmw) SyncInputFeedLine() (err error) {
	defer func(begin time.Time) {

		mw.logger.Print(
			"method", "SyncInputFeedLine",
			"input", 0,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.IFluService.SyncInputFeedLine()
	return
}

func (mw logmw) GetFeedLineUnit(fluId uuid.UUID) (flu models.FeedLineUnit, err error) {
	defer func(begin time.Time) {

		mw.logger.Print(
			"method", "SyncInputFeedLine",
			"input", fluId,
			"output", flu,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	flu, err = mw.IFluService.GetFeedLineUnit(fluId)
	return
}
