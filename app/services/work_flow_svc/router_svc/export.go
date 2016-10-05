package step_router

func newStdStepRouter() stepRouter {

	router := newStepRouter(25)

	// Connect to all the steps (in this case steps like crowd sourcing)
	router.connectAll()

	// Start the router
	start(&router)

	return router
}

var StdStepRouter = newStdStepRouter()
