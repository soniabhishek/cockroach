package step_router

import (
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/feed_line"
)

func newStdStepRouter() stepRouter {

	sr := stepRouter{
		// Bigger feedLine since all the step servers
		// pushes flu to this one only
		InQ:           feed_line.NewBig(),
		ProcessedFluQ: feed_line.New(),
	}

	// Connect to all the servers (in this case steps like crowd sourcing)
	sr.connectAll()

	// Start the router
	sr.start()

	return sr
}

var StdStepRouter = newStdStepRouter()
