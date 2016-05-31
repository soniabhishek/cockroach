package step_router

func newStdStepRouter() stepRouter {

	router := newStepRouter()

	// Connect to all the servers (in this case steps like crowd sourcing)
	router.connectAll()

	// Start the router
	router.start()

	return router
}

var StdStepRouter = newStdStepRouter()
