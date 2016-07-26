package feed_line_repo

import (
	"time"

	"github.com/crowdflux/angel/app/DAL/repositories"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"github.com/lib/pq"
)

type feedLineLogger struct {
	Db repositories.IDatabase
}

func (l *feedLineLogger) Log(fLog []models.FeedLineLog) {
	now := time.Now()
	var fLogArr []interface{} = make([]interface{}, len(fLog))

	for i, _ := range fLog {
		fLog[i].CreatedAt = pq.NullTime{now, true}
		fLogArr[i] = &fLog[i]
	}
	err := l.Db.Insert(fLogArr...)
	plog.Error("Log Error: ", err)
}
