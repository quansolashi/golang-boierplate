package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quansolashi/golang-boierplate/backend/internal/util"
	"github.com/quansolashi/golang-boierplate/backend/pkg/rabbitmq"
)

func (c *controller) queueRoutes(rg *gin.RouterGroup) {
	rg.Use(c.authentication)
	rg.POST("/publish", c.publishMessage) // `queue/publish?message=xxx``
}

func (c *controller) publishMessage(ctx *gin.Context) {
	message := util.GetQuery(ctx, "message", "")
	if message == "" {
		ctx.Status(http.StatusOK)
		return
	}

	err := c.queue.Publish(ctx, &rabbitmq.PublishParams{
		QueueName: "queue",
		Message:   message,
	})
	if err != nil {
		c.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}
