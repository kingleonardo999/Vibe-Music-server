package router

import (
	"github.com/gin-gonic/gin"
	"vibe-music-server/internal/controller"
	"vibe-music-server/internal/middleware"
)

func registerBannerRouter(r *gin.Engine, ctrl *controller.BannerCtrl) {
	g := r.Group("/admin")
	p := r.Group("/banner")
	// 组级中间件
	g.Use(middleware.AdminAuthMiddleware())
	{
		g.POST("/getAllBanners", ctrl.GetAllBanners)
		g.POST("/addBanner", ctrl.AddBanner)
		g.PATCH("/updateBanner/:id", ctrl.UpdateBanner)
		g.PATCH("/updateBannerStatus/:id", ctrl.UpdateBannerStatus)
		g.DELETE("/deleteBanner/:id", ctrl.DeleteBanner)
		g.DELETE("/deleteBanners", ctrl.DeleteBanners)
	}
	{
		p.GET("getBannerList", ctrl.GetBannerList)
	}
}
