package auther

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ShowAuthenticationErrorOverHttp(c *gin.Context, msg string) {
	c.Header("Authenication error", msg)
	c.AbortWithStatus(http.StatusUnauthorized)
}

func ShowInternalErrorOverHttp(c *gin.Context, msg string) {
	c.Header("Internal Auth error", msg)
	c.AbortWithStatus(http.StatusInternalServerError)
}
