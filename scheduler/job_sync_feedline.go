package scheduler

import "github.com/crowdflux/angel/app/services/flu_svc"

const syncFeedLineCron = "0/20 * * * * *"

func jobSyncFeedLine() {
	f := flu_svc.New()
	f.SyncInputFeedLine()
}
