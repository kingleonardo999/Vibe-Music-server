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

type PlaylistCtrl struct {
	playlistService *service.PlaylistService
}

func NewPlaylistCtrl(playlistService *service.PlaylistService) *PlaylistCtrl {
	return &PlaylistCtrl{playlistService: playlistService}
}

func (p *PlaylistCtrl) GetAllPlaylists(c *gin.Context) {
	var playlistDTO dto.PlaylistDTO
	if err := c.ShouldBindJSON(&playlistDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, p.playlistService.GetAllPlaylists(&playlistDTO))
}

func (p *PlaylistCtrl) GetRecommendedPlaylists(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusOK, p.playlistService.GetRecommendedPlaylists(nil))
		return
	}
	c.JSON(http.StatusOK, p.playlistService.GetRecommendedPlaylists(claims.(*util.Claims)))
}

func (p *PlaylistCtrl) GetPlaylistDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	playlistId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusOK, p.playlistService.GetPlaylistDetail(playlistId, nil))
		return
	}
	c.JSON(http.StatusOK, p.playlistService.GetPlaylistDetail(playlistId, claims.(*util.Claims)))
}
