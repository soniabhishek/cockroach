package api

import (
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent"
)

func NewRelicMiddleware(app newrelic.Application) gin.HandlerFunc {
	return func(c *gin.Context) {

		txn := app.StartTransaction(c.HandlerName(), c.Writer, c.Request)
		defer txn.End()

		c.Next()
	}
}
