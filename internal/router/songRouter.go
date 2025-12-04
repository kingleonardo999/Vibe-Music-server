package router

import (
	"github.com/gin-gonic/gin"
	"vibe-music-server/internal/controller"
)

func registerSongRouter(r *gin.Engine, ctrl *controller.SongCtrl) {
	g := r.Group("/song")
	{
		g.POST("/getAllSongs", ctrl.GetAllSongs)
		g.GET("/getRecommendedSongs", ctrl.GetRecommendedSongs)
		g.GET("/getSongDetail/:id", ctrl.GetSongDetail)
	}
}
