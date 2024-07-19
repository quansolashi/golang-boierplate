package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/quansolashi/golang-boierplate/backend/pkg/cors"
)

func (a *app) newRouter() *gin.Engine {
	rt := gin.New()
	rt.Use(cors.NewGinMiddleware())
	a.web.Routes(rt.Group("/api"))
	return rt
}
