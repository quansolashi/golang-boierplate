package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *controller) heartbeatRoutes(rg *gin.RouterGroup) {
	rg.GET("", c.index)
}

func (c *controller) index(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "health",
	})
}
