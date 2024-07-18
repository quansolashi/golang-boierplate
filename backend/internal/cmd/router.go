package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/quansolashi/message-extractor/backend/pkg/cors"
)

func (a *app) newRouter() *gin.Engine {
	rt := gin.New()
	rt.Use(cors.NewGinMiddleware())
	a.web.Routes(rt.Group("/api"))
	return rt
}
