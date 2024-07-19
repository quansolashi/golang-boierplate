package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quansolashi/message-extractor/backend/internal/util"
)

func (c *controller) userRoutes(rg *gin.RouterGroup) {
	rg.GET("", c.userIndex)
	rg.GET("/:userId", c.showUser)
}

func (c *controller) userIndex(ctx *gin.Context) {
	users, err := c.db.User.List(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

func (c *controller) showUser(ctx *gin.Context) {
	userID, err := util.GetParamUint64(ctx, "userId")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
	}

	user, err := c.db.User.Get(ctx, userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}
