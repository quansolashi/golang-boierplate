package cors

import (
	"github.com/gin-gonic/gin"
	wrapper "github.com/rs/cors/wrapper/gin"
)

func NewGinMiddleware() gin.HandlerFunc {
	// from cors options to gin HandlerFunc
	options := NewOptions()
	return wrapper.New(options)
}
