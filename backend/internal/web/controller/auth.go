package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/quansolashi/golang-boierplate/backend/internal/database"
	"github.com/quansolashi/golang-boierplate/backend/internal/entity"
	"github.com/quansolashi/golang-boierplate/backend/internal/web/request"
	"github.com/quansolashi/golang-boierplate/backend/internal/web/response"
	"github.com/quansolashi/golang-boierplate/backend/pkg/security"
)

func (c *controller) authRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", c.login)
}

func (c *controller) login(ctx *gin.Context) {
	req := &request.LoginRequest{}
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		c.badRequest(ctx, err)
		return
	}

	user, err := c.db.User.GetByEmail(ctx, req.Email)
	if errors.Is(err, database.ErrNotFound) {
		c.unauthorized(ctx, err)
		return
	}
	if err != nil {
		c.httpError(ctx, err)
		return
	}

	if err = security.VerifyPassword(user.Password, req.Password); err != nil {
		c.unauthorized(ctx, err)
		return
	}

	res, err := c.newUserInfo(user)
	if err != nil {
		c.httpError(ctx, err)
		return
	}
	ctx.JSON(200, gin.H{
		"data": res,
	})
}

func (c *controller) newUserInfo(user *entity.User) (*response.LoginResponse, error) {
	token, err := c.auth.CreateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	res := &response.LoginResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: token,
	}
	return res, nil
}
