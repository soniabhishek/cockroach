package scheduler

import "gitlab.com/playment-main/support/app/services/flu_svc"

const syncFeedLineCron = "0/20 * * * * *"

func jobSyncFeedLine() {
	f := flu_svc.New()
	f.SyncInputFeedLine()
}
