package controller

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/providers/google"
	"github.com/quansolashi/golang-boierplate/backend/internal/database"
	"github.com/quansolashi/golang-boierplate/backend/internal/entity"
	"github.com/quansolashi/golang-boierplate/backend/internal/util"
	"github.com/quansolashi/golang-boierplate/backend/pkg/auth"
	"github.com/quansolashi/golang-boierplate/backend/pkg/rabbitmq"
	"github.com/quansolashi/golang-boierplate/backend/pkg/redis"
)

const (
	loginWithGoogleCallbackURL string = "/callback/google"
)

type Controller interface {
	Routes(ctx *gin.RouterGroup)
}

type Params struct {
	DB               *database.Database
	Redis            *redis.Client
	RabbitMQ         rabbitmq.Client
	LocalTokenSecret string
	GoogleAPIKey     string
	GoogleAPISecret  string
	WebURL           string
}

type controller struct {
	db     *database.Database
	redis  *redis.Client
	auth   auth.LocalClient
	google *google.Provider
	queue  rabbitmq.Client
}

func NewController(params *Params) Controller {
	return &controller{
		db:    params.DB,
		redis: params.Redis,
		auth:  auth.NewLocalClient(params.LocalTokenSecret),
		google: google.New(
			params.GoogleAPIKey,
			params.GoogleAPISecret,
			fmt.Sprintf("%s%s", params.WebURL, loginWithGoogleCallbackURL),
			"email",
			"profile",
		),
		queue: params.RabbitMQ,
	}
}

func (c *controller) Routes(rg *gin.RouterGroup) {
	c.heartbeatRoutes(rg.Group("/heartbeat"))

	c.authRoutes(rg.Group("/auth"))
	c.userRoutes(rg.Group("/users"))
	c.queueRoutes(rg.Group("/queues"))
}

func (c *controller) authentication(ctx *gin.Context) {
	token := ctx.Request.Header.Get("Authorization")
	if token == "" {
		c.unauthorized(ctx, errors.New("auth: token is not found"))
		return
	}
	tokenString := strings.Split(token, "Bearer ")
	if len(tokenString) < 2 {
		c.unauthorized(ctx, errors.New("auth: token is invalid"))
		return
	}

	claims, err := c.auth.VerifyToken(tokenString[1])
	if err != nil {
		c.unauthorized(ctx, err)
		return
	}
	ctx.Request.Header.Set("userId", strconv.FormatUint(claims.ID, 10))
	ctx.Next()
}

//nolint:unused // use later
func (c *controller) getUserID(ctx *gin.Context) (uint64, error) {
	userID := ctx.GetHeader("userId")
	if userID == "" {
		return 0, nil
	}
	return strconv.ParseUint(userID, 10, 64)
}

//nolint:unused // use later
func (c *controller) currentUser(ctx *gin.Context) (*entity.User, error) {
	userID, err := c.getUserID(ctx)
	if err != nil {
		return nil, err
	}
	user, err := c.db.User.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (c *controller) httpError(ctx *gin.Context, err error) {
	res, code := util.NewErrorResponse(err)
	ctx.AbortWithStatusJSON(code, res)
}

func (c *controller) unauthorized(ctx *gin.Context, err error) {
	c.httpError(ctx, fmt.Errorf("%w: %s", util.ErrUnauthenticated, err.Error()))
}

func (c *controller) badRequest(ctx *gin.Context, err error) {
	c.httpError(ctx, fmt.Errorf("%w: %s", util.ErrInvalidArgument, err.Error()))
}
