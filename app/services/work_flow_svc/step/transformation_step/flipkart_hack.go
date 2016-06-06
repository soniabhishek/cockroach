package transformation_step

import (
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/feed_line"
)

var approveJson = models.JsonFake{
	"action":  "reject",
	"approve": "approve",
}

func flipkartHack(flu feed_line.FLU) {
	flu.Build.Merge(models.JsonFake{
		"result": approveJson,
	})
	StdTransformationStep.finishFlu(flu)
}
