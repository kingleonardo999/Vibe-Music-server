package router

import (
	"github.com/gin-gonic/gin"
	"vibe-music-server/internal/controller"
	"vibe-music-server/internal/middleware"
)

func registerFeedbackRouter(r *gin.Engine, ctrl *controller.FeedbackCtrl) {
	g := r.Group("/admin")
	p := r.Group("/feedback")
	{
		g.POST("/getAllFeedbacks", ctrl.GetAllFeedbacks)
		g.DELETE("/deleteFeedback/:id", ctrl.DeleteFeedback)
		g.DELETE("/deleteFeedbacks", ctrl.DeleteFeedbacks)
	}
	p.Use(middleware.AuthMiddleware())
	{
		p.POST("/addFeedback", ctrl.AddFeedback)
	}
}
