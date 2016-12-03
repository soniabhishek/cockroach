package transformation_step_svc

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models"
)

var approveJson = models.JsonF{
	"action":  "approve",
	"approve": "approve",
}

func flipkartHack(flu feed_line.FLU) {
	flu.Build.Merge(models.JsonF{
		"result": approveJson,
	})
	StdTransformationStep.finishFlu(flu)
}
