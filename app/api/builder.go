package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/playment-main/angel/app/api/handlers"
	"gitlab.com/playment-main/angel/app/config"
	"gitlab.com/playment-main/angel/app/services/flu_svc/flu_svc_transport"
)

func Build() {

	if config.GetEnvironment() == config.Production {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Creates a router without any middleware by default
	r := gin.New()

	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	//Api prefix
	api := r.Group("/api/v0")

	api.Use()
	{

		api.POST("/bulkdownloadimages", handlers.BulkDownloadImages)
		api.POST("/bulkdownloadimagesfromcsv", handlers.BulkDownloadedImagesFromCSV)

		flu_svc_transport.AddHttpTransport(api)
	}

	r.Run(":8999") // listen and serve on 127.0.0.1:8999

}
