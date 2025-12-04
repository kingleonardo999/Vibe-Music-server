package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"vibe-music-server/internal/pkg/result"
	"vibe-music-server/internal/pkg/result/consts"
	"vibe-music-server/internal/pkg/util"
)

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, result.Error[result.Nil](consts.NotLogin))
			c.Abort()
			return
		}
		adminClaims, ok := claims.(*util.Claims)
		if !ok || adminClaims == nil || adminClaims.Role != consts.AdminRole {
			c.JSON(http.StatusForbidden, result.Error[result.Nil](consts.NoPermission))
			c.Abort()
			return
		}
		c.Next()
	}
}
