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

type SongCtrl struct {
	songService *service.SongService
}

func NewSongCtrl(songService *service.SongService) *SongCtrl {
	return &SongCtrl{
		songService: songService,
	}
}

func (s *SongCtrl) GetAllSongs(c *gin.Context) {
	var songDTO dto.SongDTO
	if err := c.ShouldBindJSON(&songDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusOK, s.songService.GetAllSongs(&songDTO, nil))
		return
	}
	c.JSON(http.StatusOK, s.songService.GetAllSongs(&songDTO, claims.(*util.Claims)))
}

func (s *SongCtrl) GetRecommendedSongs(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusOK, s.songService.GetRecommendedSongs(nil))
		return
	}
	c.JSON(http.StatusOK, s.songService.GetRecommendedSongs(claims.(*util.Claims)))
}

func (s *SongCtrl) GetSongDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	songId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusOK, s.songService.GetSongDetail(songId, nil))
		return
	}
	c.JSON(http.StatusOK, s.songService.GetSongDetail(songId, claims.(*util.Claims)))
}
