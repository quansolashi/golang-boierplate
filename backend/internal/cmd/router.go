package cmd

import (
	ginzip "github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/quansolashi/golang-boierplate/backend/pkg/cors"
)

func (a *app) newRouter() *gin.Engine {
	rt := gin.New()
	rt.Use(cors.NewGinMiddleware())
	rt.Use(ginzip.Gzip(ginzip.DefaultCompression))

	a.web.Routes(rt.Group("/api"))
	return rt
}
