package middleware

import (
	"strings"

	"pet-social-game/server/internal/pkg/auth"
	"pet-social-game/server/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

const CtxUserID = "userId"
const CtxPlatform = "platform"

func JWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" {
			response.Unauthorized(c, "missing authorization")
			c.Abort()
			return
		}
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Unauthorized(c, "invalid authorization header")
			c.Abort()
			return
		}
		claims, err := auth.ParseToken(secret, parts[1])
		if err != nil {
			response.Unauthorized(c, "invalid token")
			c.Abort()
			return
		}
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxPlatform, claims.Platform)
		c.Next()
	}
}

func UserID(c *gin.Context) uint64 {
	v, _ := c.Get(CtxUserID)
	id, _ := v.(uint64)
	return id
}
