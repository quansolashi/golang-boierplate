package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/quansolashi/golang-boierplate/backend/internal/database"
)

type Controller interface {
	Routes(ctx *gin.RouterGroup)
}

type Params struct {
	DB *database.Database
}

type controller struct {
	db *database.Database
}

func NewController(params *Params) Controller {
	return &controller{
		db: params.DB,
	}
}

func (c *controller) Routes(rg *gin.RouterGroup) {
	c.heartbeatRoutes(rg.Group("/heartbeat"))

	c.userRoutes(rg.Group("/users"))
}
