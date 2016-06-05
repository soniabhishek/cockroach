package step_router

func newStdStepRouter() stepRouter {

	router := newStepRouter(10)

	// Connect to all the steps (in this case steps like crowd sourcing)
	router.connectAll()

	// Start the router
	router.start()

	return router
}

var StdStepRouter = newStdStepRouter()
