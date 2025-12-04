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

type FavoriteCtrl struct {
	favoriteService *service.FavoriteService
}

func NewFavoriteCtrl(favoriteService *service.FavoriteService) *FavoriteCtrl {
	return &FavoriteCtrl{
		favoriteService: favoriteService,
	}
}

// CollectSong 收藏歌曲
// need authMiddleware
func (f *FavoriteCtrl) CollectSong(c *gin.Context) {
	songIdStr := c.Query("songId")
	if songIdStr == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	songId, err := strconv.ParseUint(songIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	claims, exist := c.Get("claims")
	if !exist {
		c.JSON(http.StatusUnauthorized, result.Error[result.Nil](consts.NotLogin))
		return
	}
	c.JSON(http.StatusOK, f.favoriteService.CollectSong(songId, claims.(*util.Claims)))
}

// CancelCollectSong 取消收藏歌曲
// need authMiddleware
func (f *FavoriteCtrl) CancelCollectSong(c *gin.Context) {
	songIdStr := c.Query("songId")
	if songIdStr == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	songId, err := strconv.ParseUint(songIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	claims, exist := c.Get("claims")
	if !exist {
		c.JSON(http.StatusUnauthorized, result.Error[result.Nil](consts.NotLogin))
		return
	}
	c.JSON(http.StatusOK, f.favoriteService.CancelCollectSong(songId, claims.(*util.Claims)))
}

// GetFavoritePlaylists 获取用户收藏的歌单
// need authMiddleware
func (f *FavoriteCtrl) GetFavoritePlaylists(c *gin.Context) {
	var playlistDTO dto.PlaylistDTO
	if err := c.ShouldBindJSON(&playlistDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	claims, exist := c.Get("claims")
	if !exist {
		c.JSON(http.StatusUnauthorized, result.Error[result.Nil](consts.NotLogin))
		return
	}
	c.JSON(http.StatusOK, f.favoriteService.GetUserFavoritePlaylists(&playlistDTO, claims.(*util.Claims)))
}

// CollectPlaylist 收藏歌单
// need authMiddleware
func (f *FavoriteCtrl) CollectPlaylist(c *gin.Context) {
	playlistIdStr := c.Query("playlistId")
	if playlistIdStr == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	playlistId, err := strconv.ParseUint(playlistIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	claims, exist := c.Get("claims")
	if !exist {
		c.JSON(http.StatusUnauthorized, result.Error[result.Nil](consts.NotLogin))
		return
	}
	c.JSON(http.StatusOK, f.favoriteService.CollectPlaylist(playlistId, claims.(*util.Claims)))
}

// CancelCollectPlaylist 取消收藏歌单
// need authMiddleware
func (f *FavoriteCtrl) CancelCollectPlaylist(c *gin.Context) {
	playlistIdStr := c.Query("playlistId")
	if playlistIdStr == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	playlistId, err := strconv.ParseUint(playlistIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	claims, exist := c.Get("claims")
	if !exist {
		c.JSON(http.StatusUnauthorized, result.Error[result.Nil](consts.NotLogin))
		return
	}
	c.JSON(http.StatusOK, f.favoriteService.CancelCollectPlaylist(playlistId, claims.(*util.Claims)))
}
