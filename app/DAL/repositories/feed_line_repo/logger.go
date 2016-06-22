package feed_line_repo

import (
	"github.com/lib/pq"
	"gitlab.com/playment-main/angel/app/DAL/repositories"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/plog"
	"time"
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
