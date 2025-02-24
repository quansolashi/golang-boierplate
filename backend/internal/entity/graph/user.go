package graph

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/quansolashi/golang-boierplate/backend/pkg/auth"
)

type contextKey struct {
	name string
}

var userIDCtxKey = &contextKey{
	name: "UID",
}

type UserID struct {
	uid uint64
}

func (u *UserID) UID() uint64 {
	return u.uid
}

func TokenAuthMiddleware(auth auth.LocalClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			return
		}
		tokenString := strings.Split(token, "Bearer ")
		if len(tokenString) < 2 {
			return
		}

		// FIXME: verify token and build value from claims
		claims, err := auth.VerifyToken(tokenString[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &gin.H{
				"error": "unauthorized",
			})
			return
		}

		uid := &UserID{
			uid: claims.ID,
		}
		ctx := context.WithValue(c.Request.Context(), userIDCtxKey, uid)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func UserIDFromContext(ctx context.Context) *UserID {
	value := ctx.Value(userIDCtxKey)
	if value != nil {
		if token, ok := value.(*UserID); ok {
			return token
		}
	}
	return nil
}
