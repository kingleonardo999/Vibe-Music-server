package router

import (
	"github.com/gin-gonic/gin"
	"vibe-music-server/internal/controller"
	"vibe-music-server/internal/middleware"
)

func registerCommentRouter(r *gin.Engine, ctrl *controller.CommentCtrl) {
	g := r.Group("/comment")
	{
		g.PATCH("/likeComment/:id", ctrl.LikeComment)
		g.PATCH("/cancelLikeComment/:id", ctrl.CancelLikeComment)
	}
	g.Use(middleware.AuthMiddleware())
	{
		g.POST("/addSongComment", ctrl.AddSongComment)
		g.POST("/addPlaylistComment", ctrl.AddPlaylistComment)
		g.DELETE("/deleteComment/:id", ctrl.DeleteComment)
	}
}
