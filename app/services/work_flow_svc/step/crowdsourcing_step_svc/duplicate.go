package crowdsourcing_step_svc

import (
	"github.com/crowdflux/angel/app/DAL/clients/postgres"
	"github.com/crowdflux/angel/app/models/uuid"
	"sync"
)

var onn sync.Once
var duplicateFluIds []uuid.UUID

func isDuplicate(fluId uuid.UUID) bool {

	onn.Do(func() {
		pgClient := postgres.GetPostgresClient()

		_, err := pgClient.Select(&duplicateFluIds, `
	SELECT DISTINCT (flu_id) from feed_line_log
	WHERE event = 3 AND message = 'crowdySendFailure';`)
		if err != nil {
			panic(err)
		}
	})

	for _, dFlu := range duplicateFluIds {
		if fluId == dFlu {
			return true
		}
	}
	return false
}
