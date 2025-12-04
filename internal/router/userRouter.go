package router

import (
	"github.com/gin-gonic/gin"
	"vibe-music-server/internal/controller"
	"vibe-music-server/internal/middleware"
)

func registerUserRouter(r *gin.Engine, ctrl *controller.UserCtrl) {
	g := r.Group("/user")
	{
		g.GET("/sendVerificationCode", ctrl.SendVerificationCode)
		g.POST("/register", ctrl.Register)
		g.POST("/login", ctrl.Login)
		g.PATCH("/resetUserPassword", ctrl.ResetUserPassword)
	}
	g.Use(middleware.AuthMiddleware())
	{
		g.GET("/getUserInfo", ctrl.GetUserInfo)
		g.PUT("/updateUserInfo", ctrl.UpdateUserInfo)
		g.PATCH("/updateUserAvatar", ctrl.UpdateUserAvatar)
		g.PATCH("/updateUserPassword", ctrl.UpdateUserPassword)
		g.POST("logout", ctrl.Logout)
		g.DELETE("/deleteAccount", ctrl.DeleteAccount)
	}
}
