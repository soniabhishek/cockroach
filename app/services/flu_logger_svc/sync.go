package flu_logger_svc

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"time"
)

func dbLogSyncer() {

	// Get log receiver channel
	feedlineLoggerChannel := feed_line.GetFeedlineLoggerChannel()
	flogs := feedlineLoggerChannel.Receiver()

	bufferedFlogs := []models.FeedLineLog{}
	flr := feed_line_repo.NewLogger()

	var saveInDb = func() {
		err := flr.Log(bufferedFlogs)

		if err != nil {
			plog.Error("Feedline logger", err, "error saving logs to db")
		} else {
			feedlineLoggerChannel.ConfirmReceivedProcessed()
			bufferedFlogs = nil
		}
	}

	ticker := time.Tick(time.Duration(2) * time.Second)

	for {

		select {

		case flog := <-flogs:

			// if channel is returning flogs
			// then store it in a temp buffer
			bufferedFlogs = append(bufferedFlogs, flog)

			if len(bufferedFlogs) > 10000 {
				saveInDb()
			}

		case <-ticker:

			// if channel stops returning flus
			// then put them in db
			if len(bufferedFlogs) > 0 {
				saveInDb()
			}

		}
	}
}

func runOnce() uint {
	go func() {

		time.Sleep(time.Duration(3) * time.Second)
		dbLogSyncer()
	}()

	return 0
}

var _ = runOnce()
