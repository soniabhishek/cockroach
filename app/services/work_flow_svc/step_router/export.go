package step_router

func newStdStepRouter() stepRouter {

	sr := new()

	// Connect to all the servers (in this case steps like crowd sourcing)
	sr.connectAll()

	// Start the router
	sr.start()

	return sr
}

var StdStepRouter = newStdStepRouter()
