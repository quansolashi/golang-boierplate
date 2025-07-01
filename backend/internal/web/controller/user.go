package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quansolashi/golang-boierplate/backend/internal/util"
	"github.com/quansolashi/golang-boierplate/backend/internal/web/service"
)

func (c *controller) userRoutes(rg *gin.RouterGroup) {
	rg.Use(c.authentication)
	rg.GET("", c.userIndex)
	rg.GET("/:userId", c.showUser)
}

// @Summary     user index
// @Description list users
// @Tags        User
// @Router      /users [get]
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Success     200 {object} response.Users
// @Failure     400 {object} util.ErrorResponse
// @Failure     404 {object} util.ErrorResponse
// @Failure     500 {object} util.ErrorResponse
func (c *controller) userIndex(ctx *gin.Context) {
	users, err := c.db.User.List(ctx)
	if err != nil {
		c.httpError(ctx, err)
	}

	res := service.NewUsers(users).Response()
	ctx.JSON(http.StatusOK, res)
}

// @Summary     show user
// @Description detail user
// @Tags        User
// @Router      /users/{userId} [get]
// @Param				userId path uint64 true "User ID"
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Success     200 {object} response.Users
// @Failure     400 {object} util.ErrorResponse
// @Failure     404 {object} util.ErrorResponse
// @Failure     500 {object} util.ErrorResponse
func (c *controller) showUser(ctx *gin.Context) {
	userID, err := util.GetParamUint64(ctx, "userId")
	if err != nil {
		c.badRequest(ctx, err)
	}

	user, err := c.db.User.Get(ctx, userID)
	if err != nil {
		c.httpError(ctx, err)
	}

	res := service.NewUser(user).Response()
	ctx.JSON(http.StatusOK, res)
}
