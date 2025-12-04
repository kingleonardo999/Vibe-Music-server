package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"vibe-music-server/internal/pkg/result/consts"
	"vibe-music-server/internal/pkg/util"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsI, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, consts.NotLogin)
			c.Abort()
			return
		}
		claims, ok := claimsI.(*util.Claims)
		if !ok || claims == nil {
			c.JSON(http.StatusUnauthorized, consts.NoPermission)
			c.Abort()
			return
		}
		c.Next()
	}
}
