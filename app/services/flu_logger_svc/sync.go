package flu_logger_svc

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"time"
)

func dbLogSyncer(fluLogger feed_line_repo.IFluLogger) {

	// Get log receiver channel
	feedlineLoggerChannel := feed_line.GetFeedlineLoggerChannel()
	flogs := feedlineLoggerChannel.Receiver()

	bufferedFlogs := []models.FeedLineLog{}

	var save = func() {
		err := fluLogger.Log(bufferedFlogs)

		if err != nil {
			plog.Error("Feedline logger", err, "error saving logs to db")
		} else {

			// confirm to the channel that the batch processing has been
			// done
			feedlineLoggerChannel.ConfirmReceivedProcessed()

			// clear the buffer for next ones
			bufferedFlogs = nil
		}
	}

	ticker := time.Tick(time.Duration(10) * time.Second)

	// the below logic basically saves the flu logs in db
	// it saves whatever is there every 10 seconds
	// or if buffered logs count reaches 10000

	for {

		select {

		case flog := <-flogs:

			// if channel is returning flogs
			// then store it in a temp buffer
			bufferedFlogs = append(bufferedFlogs, flog)

			if len(bufferedFlogs) >= 10000 {
				save()
			}

		case <-ticker:

			// if channel stops returning flus
			// then put them in db
			if len(bufferedFlogs) > 0 {
				save()
			}

		}
	}
}

func runOnceOnAppStart() uint {
	go func() {

		// delayed start just for convenience in logs
		time.Sleep(time.Duration(3) * time.Second)
		dbLogSyncer(feed_line_repo.NewLogger())
	}()

	return 0
}

var _ = runOnceOnAppStart()

func mockDbLogSyncer(fluLogger feed_line_repo.IFluLogger) {

	go func() { dbLogSyncer(fluLogger) }()
}
