package step_router

func newStdStepRouter() stepRouter {

	router := newStepRouter(50)

	// Connect to all the steps (in this case steps like crowd sourcing)
	router.connectAll()

	// Start the router
	start(&router)

	return router
}

var StdStepRouter = newStdStepRouter()
