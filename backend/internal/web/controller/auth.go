package controller

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/quansolashi/golang-boierplate/backend/internal/database"
	"github.com/quansolashi/golang-boierplate/backend/internal/entity"
	"github.com/quansolashi/golang-boierplate/backend/internal/web/request"
	"github.com/quansolashi/golang-boierplate/backend/internal/web/response"
	"github.com/quansolashi/golang-boierplate/backend/pkg/security"
)

func (c *controller) authRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", c.login)
	rg.GET("/google", c.loginWithGoogle)
	rg.POST("/google/callback", c.loginWithGoogleCallback)
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

	//temporary: write login time to redis
	userIDStr := strconv.FormatUint(user.ID, 10)
	layout := "2006-01-02 15:04:05"
	err = c.redis.Set(ctx, userIDStr, time.Now().Format(layout))
	if err != nil {
		c.httpError(ctx, err)
		return
	}

	res, err := c.newUserInfo(ctx, user)
	if err != nil {
		c.httpError(ctx, err)
		return
	}
	ctx.JSON(200, gin.H{
		"data": res,
	})
}

func (c *controller) loginWithGoogle(ctx *gin.Context) {
	ctx.Request = c.getContextWithGoogle(ctx)
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

func (c *controller) loginWithGoogleCallback(ctx *gin.Context) {
	ctx.Request = c.getContextWithGoogle(ctx)

	auth, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		c.unauthorized(ctx, err)
		return
	}
	user, err := c.db.User.GetByEmail(ctx, auth.Email)
	if err != nil {
		c.httpError(ctx, err)
		return
	}

	res, err := c.newUserInfo(ctx, user)
	if err != nil {
		c.httpError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *controller) getContextWithGoogle(ctx *gin.Context) *http.Request {
	const key, value string = "provider", "google"
	goth.UseProviders(c.google)
	//nolint:revive,staticcheck // variables must be string
	return ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), key, value))
}

func (c *controller) newUserInfo(ctx context.Context, user *entity.User) (*response.LoginResponse, error) {
	token, err := c.auth.CreateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	//temporary: get last accessed at of user from redis
	userIDStr := strconv.FormatUint(user.ID, 10)
	accessedAtStr, err := c.redis.Get(ctx, userIDStr)
	if err != nil {
		return nil, err
	}
	layout := "2006-01-02 15:04:05"
	accessedAt, err := time.Parse(layout, accessedAtStr.(string))
	if err != nil {
		return nil, err
	}

	res := &response.LoginResponse{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		AccessToken:    token,
		LastAccessedAt: accessedAt,
	}
	return res, nil
}
