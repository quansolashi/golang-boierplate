package controller

import "github.com/gin-gonic/gin"

type Controller interface {
	Routes(ctx *gin.RouterGroup)
}

type Params struct {
}

type controller struct{}

func NewController(params *Params) Controller {
	return &controller{}
}

func (c *controller) Routes(rg *gin.RouterGroup) {
	c.heartbeatRoutes(rg.Group("/heartbeat"))
}
