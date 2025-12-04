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

type CommentCtrl struct {
	commentService *service.CommentService
}

func NewCommentCtrl(commentService *service.CommentService) *CommentCtrl {
	return &CommentCtrl{commentService: commentService}
}

// AddSongComment 添加歌曲评论
// need authMiddleware
func (m *CommentCtrl) AddSongComment(c *gin.Context) {
	var commentSongDTO dto.CommentSongDTO
	if err := c.ShouldBindJSON(&commentSongDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	claims, exist := c.Get("claims")
	if !exist {
		c.JSON(http.StatusUnauthorized, result.Error[result.Nil](consts.NotLogin))
		return
	}
	c.JSON(http.StatusOK, m.commentService.AddSongComment(&commentSongDTO, claims.(*util.Claims)))
}

// AddPlaylistComment 添加歌单评论
// need authMiddleware
func (m *CommentCtrl) AddPlaylistComment(c *gin.Context) {
	var commentPlaylistDTO dto.CommentPlaylistDTO
	if err := c.ShouldBindJSON(&commentPlaylistDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	claims, exist := c.Get("claims")
	if !exist {
		c.JSON(http.StatusUnauthorized, result.Error[result.Nil](consts.NotLogin))
		return
	}
	c.JSON(http.StatusOK, m.commentService.AddPlaylistComment(&commentPlaylistDTO, claims.(*util.Claims)))
}

func (m *CommentCtrl) LikeComment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	commentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, m.commentService.LikeComment(commentID))
}

func (m *CommentCtrl) CancelLikeComment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	commentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, m.commentService.CancelLikeComment(commentID))
}

// DeleteComment 删除评论
// need authMiddleware
func (m *CommentCtrl) DeleteComment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	commentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	claims, exist := c.Get("claims")
	if !exist {
		c.JSON(http.StatusUnauthorized, result.Error[result.Nil](consts.NotLogin))
		return
	}
	c.JSON(http.StatusOK, m.commentService.DeleteComment(commentID, claims.(*util.Claims)))
}
