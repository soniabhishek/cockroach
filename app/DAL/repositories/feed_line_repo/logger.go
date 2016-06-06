package feed_line_repo

import (
	"gitlab.com/playment-main/angel/app/DAL/repositories"
	"gitlab.com/playment-main/angel/app/models"
	"time"
)

type feedLineLogger struct {
	Db repositories.IDatabase
}

func (l *feedLineLogger) Log(fLog models.FeedLineLog) {
	fLog.CreatedAt.Time = time.Now()
	go l.Db.Insert(&fLog)
}
