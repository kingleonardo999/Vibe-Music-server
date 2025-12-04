package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"vibe-music-server/internal/model/dto"
	"vibe-music-server/internal/pkg/result"
	"vibe-music-server/internal/pkg/result/consts"
	"vibe-music-server/internal/pkg/util"
	"vibe-music-server/internal/service"
)

type ArtistCtrl struct {
	artistService *service.ArtistService
}

func NewArtistCtrl(artistService *service.ArtistService) *ArtistCtrl {
	return &ArtistCtrl{
		artistService: artistService,
	}
}

func (a *ArtistCtrl) GetAllArtists(c *gin.Context) {
	var artistDTO dto.ArtistDTO
	if err := c.ShouldBindJSON(&artistDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.artistService.GetAllArtists(&artistDTO))
}

func (a *ArtistCtrl) GetRandomArtists(c *gin.Context) {
	c.JSON(http.StatusOK, a.artistService.GetRandomArtists())
}

// GetArtistDetail
// need authMiddleware
func (a *ArtistCtrl) GetArtistDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	artistID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, a.artistService.GetArtistDetail(artistID, nil))
		return
	}
	c.JSON(http.StatusOK, a.artistService.GetArtistDetail(artistID, claims.(*util.Claims)))
}
