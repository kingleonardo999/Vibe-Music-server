package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"vibe-music-server/internal/model/dto"
	"vibe-music-server/internal/model/entity"
	"vibe-music-server/internal/pkg/result"
	"vibe-music-server/internal/pkg/result/consts"
	"vibe-music-server/internal/service"
)

type BannerCtrl struct {
	bannerService *service.BannerService
	minioService  *service.MinioService
}

func NewBannerCtrl(bannerService *service.BannerService, minioService *service.MinioService) *BannerCtrl {
	return &BannerCtrl{
		bannerService: bannerService,
		minioService:  minioService,
	}
}

func (b *BannerCtrl) GetAllBanners(c *gin.Context) {
	var bannerDTO dto.BannerDTO
	if err := c.ShouldBindQuery(&bannerDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, b.bannerService.GetAllBanners(&bannerDTO))
}

func (b *BannerCtrl) AddBanner(c *gin.Context) {
	banner, err := c.FormFile("banner")
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	bannerUrl, err := b.minioService.UploadFile(banner, "banners")
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.Error[result.Nil](consts.InternalError))
		return
	}
	c.JSON(http.StatusOK, b.bannerService.AddBanner(bannerUrl))
}

func (b *BannerCtrl) UpdateBanner(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	bannerId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	banner, err := c.FormFile("banner")
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	bannerUrl, err := b.minioService.UploadFile(banner, "banners")
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.Error[result.Nil](consts.InternalError))
		return
	}
	c.JSON(http.StatusOK, b.bannerService.UpdateBanner(bannerId, bannerUrl))
}

func (b *BannerCtrl) UpdateBannerStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	bannerId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	statusStr := c.Request.FormValue("status")
	status, err := strconv.Atoi(statusStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, b.bannerService.UpdateBannerStatus(bannerId, entity.BannerStatus(status)))
}

func (b *BannerCtrl) DeleteBanner(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	bannerId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, b.bannerService.DeleteBanner(bannerId))
}

func (b *BannerCtrl) DeleteBanners(c *gin.Context) {
	var ids []uint64
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, b.bannerService.DeleteBanners(ids))
}

func (b *BannerCtrl) GetBannerList(c *gin.Context) {
	c.JSON(http.StatusOK, b.bannerService.GetBannerList())
}
