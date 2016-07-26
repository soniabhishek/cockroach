package transformation_step

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/services/work_flow_svc/feed_line"
)

var approveJson = models.JsonFake{
	"action":  "approve",
	"approve": "approve",
}

func flipkartHack(flu feed_line.FLU) {
	flu.Build.Merge(models.JsonFake{
		"result": approveJson,
	})
	StdTransformationStep.finishFlu(flu)
}
