package crowdsourcing_step_svc

import (
	"github.com/crowdflux/angel/app/DAL/clients/postgres"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"time"
)

func runMe() bool {

	go func() {
		time.Sleep(time.Duration(10) * time.Second)

		pgClient := postgres.GetPostgresClient()

		var flus []models.FeedLineUnit

		plog.Info("fucked", "fuckup started")

		_, err := pgClient.Select(&flus, `SELECT *
FROM feed_line WHERE id
  IN (select
        DISTINCT fl.id
      from crowdsourcing_flu_buffer cfb
        inner join questions q
          on q.id = cfb.question_id
        inner join feed_line fl
          on fl.id = cfb.flu_id
        join micro_task_question_associators mtqa
          on mtqa.question_id = q.id
        join crowdsourcing_step_configuration csc
          on csc.step_id = fl.step_id
      where
        fl.created_at > '2016-10-01 00:00:00+05:30'
        --and fl.created_at < '2016-10-27'
        --cfb.created_at  < '2016-10-27' and  cfb.created_at > '2016-10-23'
        and cfb.is_deleted = true and q.is_active = false
        and mtqa.micro_task_id = csc.micro_task_id)
;`)
		if err != nil {
			panic(err)
		}

		plog.Info("fucked", len(flus), "started to finish")

		for _, flu := range flus {
			StdCrowdSourcingStep.finishFlu(feed_line.FLU{FeedLineUnit: flu})

		}

	}()

	return true
}

var _ = runMe()
