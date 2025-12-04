package router

import (
	"github.com/gin-gonic/gin"
	"vibe-music-server/internal/controller"
	"vibe-music-server/internal/middleware"
)

func registerArtistRouter(r *gin.Engine, ctrl *controller.ArtistCtrl) {
	g := r.Group("/artist")

	// 组级中间件
	{
		g.POST("/getAllArtists", ctrl.GetAllArtists)
		g.GET("/getRandomArtists", ctrl.GetRandomArtists)
	}
	g.Use(middleware.AdminAuthMiddleware())
	{
		g.GET("/getArtistDetail/:id", ctrl.GetArtistDetail)
	}
}
