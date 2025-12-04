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

type FeedbackCtrl struct {
	feedbackService *service.FeedbackService
}

func NewFeedbackCtrl(feedbackService *service.FeedbackService) *FeedbackCtrl {
	return &FeedbackCtrl{feedbackService: feedbackService}
}

func (f *FeedbackCtrl) GetAllFeedbacks(c *gin.Context) {
	var feedbackDTO dto.FeedbackDTO
	if err := c.ShouldBindJSON(&feedbackDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, f.feedbackService.GetAllFeedbacks(&feedbackDTO))
}

func (f *FeedbackCtrl) DeleteFeedback(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	feedbackId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, f.feedbackService.DeleteFeedback(feedbackId))
}

func (f *FeedbackCtrl) DeleteFeedbacks(c *gin.Context) {
	var ids []uint64
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, f.feedbackService.DeleteFeedbacks(ids))
}

// AddFeedback 添加反馈
// need authMiddleware
func (f *FeedbackCtrl) AddFeedback(c *gin.Context) {
	content := c.Request.FormValue("content")
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, result.Error[result.Nil](consts.NotLogin))
		return
	}
	c.JSON(http.StatusOK, f.feedbackService.AddFeedback(content, claims.(*util.Claims)))
}
