package router

import (
	"github.com/gin-gonic/gin"
	"vibe-music-server/internal/controller"
)

func registerPlaylistRouter(r *gin.Engine, ctrl *controller.PlaylistCtrl) {
	g := r.Group("/playlist")
	{
		g.POST("/getAllPlaylists", ctrl.GetAllPlaylists)
		g.GET("/getRecommendedPlaylists", ctrl.GetRecommendedPlaylists)
		g.GET("/getPlaylistDetail/:id", ctrl.GetPlaylistDetail)
	}
}
