package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/playment-main/angel/app/api/auther"
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
	{

		api.POST("/bulkdownloadimages", handlers.BulkDownloadImages)
		api.POST("/bulkdownloadimagesfromcsv", handlers.BulkDownloadedImagesFromCSV)
	}

	authorized := r.Group("/api/v0", auther.GinAuther())
	{
		flu_svc_transport.AddHttpTransport(authorized)
	}

	r.Run(":8999") // listen and serve on 127.0.0.1:8999

}
