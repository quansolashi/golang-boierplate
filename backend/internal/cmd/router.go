package cmd

import (
	ginzip "github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/quansolashi/golang-boierplate/backend/pkg/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (a *app) newRouter() *gin.Engine {
	rt := gin.New()
	rt.Use(cors.NewGinMiddleware())
	rt.Use(ginzip.Gzip(ginzip.DefaultCompression))

	// for swagger
	rt.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	a.web.Routes(rt.Group("/api"))
	a.graph.Handler(rt.Group("/graphql"))
	return rt
}
