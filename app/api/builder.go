package api

import (
	"net/http"

	"github.com/crowdflux/angel/app/api/auther"
	"github.com/crowdflux/angel/app/api/handlers"
	"github.com/crowdflux/angel/app/config"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_svc_transport"
	"github.com/crowdflux/angel/utilities/clients/api"
	"github.com/gin-gonic/gin"

	"time"

	"github.com/crowdflux/angel/app/services/work_flow_executor_svc/flu_migration_svc"
	"github.com/crowdflux/angel/app/services/work_flow_executor_svc/step/crowdsourcing_step_svc"
	"github.com/crowdflux/angel/app/services/work_flow_executor_svc/step/manual_step_svc"
	"github.com/crowdflux/angel/app/services/work_flow_svc"
	"github.com/crowdflux/angel/app/support"
	"github.com/itsjamie/gin-cors"
	"github.com/newrelic/go-agent"
)

func Build() {

	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)

	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Creates a router without any middleware by default
	r := gin.New()

	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// added cors because of int.playment.in
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	if config.IsProduction() {

		newrelicConfig := newrelic.NewConfig("Angel Server", config.NEW_RELIC_KEY.Get())
		newrelicApp, err := newrelic.NewApplication(newrelicConfig)
		if err != nil {
			panic(err)
		}

		r.Use(NewRelicMiddleware(newrelicApp))
	}

	r.StaticFS("/downloads", http.Dir(support.GetExposedDir()))

	//Api prefix
	api := r.Group("/api/v0")
	{
		api.GET("", func(c *gin.Context) {})
		api.POST("/bulkdownloadimages", handlers.BulkDownloadImages)
		api.POST("/bulkdownloadimagesfromcsv", handlers.BulkDownloadedImagesFromCSV)

		crowdsourcing_step_svc.AddHttpTransport(api)
		manual_step_svc.AddHttpTransport(api)
		utils_api.AddHttpTransport(api)
		flu_migration_svc.AddHttpTransport(api)
	}

	authorized_header := r.Group("/api/v0", auther.AuthorizeHeader())
	{
		work_flow_svc.AddHttpTransport(authorized_header)
	}

	authorized := r.Group("/api/v0", auther.GinAuther())
	{
		flu_svc_transport.AddHttpTransport(authorized)
	}

	r.Run(":8999") // listen and serve on 127.0.0.1:8999
}
