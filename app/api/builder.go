package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/playment-main/angel/app/api/auther"
	"gitlab.com/playment-main/angel/app/api/handlers"
	"gitlab.com/playment-main/angel/app/config"
	"gitlab.com/playment-main/angel/app/services/flu_svc/flu_svc_transport"
	"gitlab.com/playment-main/angel/app/services/image_svc1"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step/crowdsourcing_step/crowdsourcing_step_transport"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step/manual_step"
	"gitlab.com/playment-main/angel/utilities/clients/api"
	"net/http"
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

	fmt.Println(config.Get(config.DOWNLOAD_PATH))
	r.StaticFS("/downloadedfiles", http.Dir(config.Get(config.DOWNLOAD_PATH)))

	//Api prefix

	api := r.Group("/api/v0")
	{

		api.POST("/bulkdownloadimages", handlers.BulkDownloadImages)
		api.POST("/bulkdownloadimagesfromcsv", handlers.BulkDownloadedImagesFromCSV)

		crowdsourcing_step_transport.AddHttpTransport(api)
		manual_step.AddHttpTransport(api)
		utils_api.AddHttpTransport(api)
	}

	authorized := r.Group("/api/v0", auther.GinAuther())
	{
		flu_svc_transport.AddHttpTransport(authorized)
	}

	var _ image_svc1.IImageService

	r.Run(":8999") // listen and serve on 127.0.0.1:8999

}
