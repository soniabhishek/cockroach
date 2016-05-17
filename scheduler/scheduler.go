package scheduler

import "github.com/robfig/cron"

func Start() {

	c := cron.New()
	c.AddFunc(syncFeedLineCron, jobSyncFeedLine)
	c.Start()
}
