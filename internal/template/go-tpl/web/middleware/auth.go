package middleware

import (
	"strings"
	"{{.ModulePath}}/infra/jwt"
	"{{.ModulePath}}/logic/shared"
	"{{.ModulePath}}/web/base"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			base.Error(c, 1000, "no token")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			base.Error(c, 1001, "invalid token format")
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			base.Error(c, 1001, "invalid token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}