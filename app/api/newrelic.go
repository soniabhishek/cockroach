package api

import (
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent"
	"strings"
)

func NewRelicMiddleware(app newrelic.Application) gin.HandlerFunc {
	return func(c *gin.Context) {

		handlerStrings := strings.Split(c.HandlerName(), "/")
		handlerName := handlerStrings[len(handlerStrings)-1]

		txn := app.StartTransaction(handlerName, c.Writer, c.Request)
		defer txn.End()

		c.Next()
	}
}
