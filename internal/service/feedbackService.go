package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
	"vibe-music-server/internal/model/dto"
	"vibe-music-server/internal/model/entity"
	"vibe-music-server/internal/pkg/result"
	"vibe-music-server/internal/pkg/result/consts"
	"vibe-music-server/internal/pkg/util"
	"vibe-music-server/internal/repo"
)

type FeedbackService struct {
	feedbackRepo *repo.FeedbackRepo
}

func NewFeedbackService(feedbackRepo *repo.FeedbackRepo) *FeedbackService {
	return &FeedbackService{
		feedbackRepo: feedbackRepo,
	}
}

func (f FeedbackService) GetAllFeedbacks(feedbackDTO *dto.FeedbackDTO) result.Result[result.PageResult[entity.Feedback]] {
	retErr := result.Error[result.PageResult[entity.Feedback]]
	retSuc := result.SuccessWithData[result.PageResult[entity.Feedback]]
	pageNum := feedbackDTO.PageNum
	pageSize := feedbackDTO.PageSize
	keyword := feedbackDTO.Keyword
	startIndex := (pageNum - 1) * pageSize
	var data result.PageResult[entity.Feedback]
	templateKey := fmt.Sprintf("feedback:getAllFeedbacks:%v-%v-%v", *keyword, startIndex, pageSize)
	if util.GetCache(templateKey, &data) {
		return retSuc(consts.Success, data)
	}
	if err := f.feedbackRepo.GetAllFeedbacks(&data, keyword, startIndex, pageSize); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return retErr(consts.DataNotFound)
		}
		return retErr(consts.InternalError)
	}
	if data.Total == 0 {
		return retErr(consts.DataNotFound)
	}
	util.SetCache(templateKey, data)
	return retSuc(consts.Success, data)
}

func (f FeedbackService) DeleteFeedback(feedbackId uint64) result.Result[result.Nil] {
	if err := f.feedbackRepo.DeleteFeedback(feedbackId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result.Error[result.Nil](consts.DataNotFound)
		}
		return result.Error[result.Nil](consts.InternalError)
	}
	util.DeleteCacheByPattern("feedback:*")
	return result.Success[result.Nil](consts.Success)
}

func (f FeedbackService) DeleteFeedbacks(feedbackIds []uint64) result.Result[result.Nil] {
	if err := f.feedbackRepo.DeleteFeedbacks(feedbackIds); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result.Error[result.Nil](consts.DataNotFound)
		}
		return result.Error[result.Nil](consts.InternalError)
	}
	util.DeleteCacheByPattern("feedback:*")
	return result.Success[result.Nil](consts.Success)
}

func (f FeedbackService) AddFeedback(content string, claims *util.Claims) result.Result[result.Nil] {
	userId := claims.UserId
	feedback := entity.Feedback{
		UserID:     userId,
		Feedback:   content,
		CreateTime: time.Now(),
	}
	if err := f.feedbackRepo.AddFeedback(&feedback); err != nil {
		return result.Error[result.Nil](consts.InternalError)
	}
	util.DeleteCacheByPattern("feedback:*")
	return result.Success[result.Nil](consts.Success)
}
