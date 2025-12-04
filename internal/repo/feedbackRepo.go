package repo

import (
	"vibe-music-server/internal/model/entity"
	"vibe-music-server/internal/pkg/db"
	"vibe-music-server/internal/pkg/result"
)

type FeedbackRepo struct{}

func NewFeedbackRepo() *FeedbackRepo {
	return &FeedbackRepo{}
}

func (f FeedbackRepo) GetAllFeedbacks(data *result.PageResult[entity.Feedback], keyword *string, startIndex, pageSize int) error {
	query := db.Get().Model(&entity.Feedback{})
	if keyword != nil {
		query = query.Where("content LIKE ?", "%"+*keyword+"%")
	}
	return query.Count(&data.Total).
		Offset(startIndex).
		Limit(pageSize).
		Order("created_at DESC").
		Scan(&data.Items).Error
}

func (f FeedbackRepo) DeleteFeedback(feedbackId uint64) error {
	return db.Get().Delete(&entity.Feedback{}, feedbackId).Error
}

func (f FeedbackRepo) DeleteFeedbacks(feedbackIds []uint64) error {
	return db.Get().Delete(&entity.Feedback{}, feedbackIds).Error
}

func (f FeedbackRepo) AddFeedback(feedback *entity.Feedback) error {
	return db.Get().Create(feedback).Error
}
