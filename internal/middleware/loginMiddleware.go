package middleware

import (
	"github.com/gin-gonic/gin"
	"vibe-music-server/internal/pkg/util"
)

func LoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		var claims *util.Claims

		if tokenStr != "" {
			var err error
			claims, err = util.ParseToken(tokenStr)
			if err != nil {
				// token 无效，视为未登录（不中断请求，仅 claims=nil）
				claims = nil
			}
			// 若 token 有效，claims 已被赋值
		}
		// 无论是否传 token，都设置 claims（可能为 nil）
		c.Set("claims", claims)
		c.Next()
	}
}
