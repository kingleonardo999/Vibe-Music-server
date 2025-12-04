package router

import (
	"github.com/gin-gonic/gin"
	"vibe-music-server/internal/controller"
	"vibe-music-server/internal/middleware"
)

func registerFavoriteRouter(r *gin.Engine, ctrl *controller.FavoriteCtrl) {
	g := r.Group("/favorite")
	g.Use(middleware.AuthMiddleware())
	// song
	{
		g.POST("/getFavoriteSongs", ctrl.GetFavoritePlaylists)
		g.POST("/collectSong", ctrl.CollectSong)
		g.DELETE("/cancelCollectSong", ctrl.CancelCollectSong)
	}
	// playlist
	{
		g.POST("/getFavoritePlaylists", ctrl.GetFavoritePlaylists)
		g.POST("/collectPlaylist", ctrl.CollectPlaylist)
		g.DELETE("/cancelCollectPlaylist", ctrl.CancelCollectPlaylist)
	}
}
