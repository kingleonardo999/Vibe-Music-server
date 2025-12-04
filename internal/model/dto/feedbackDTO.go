package dto

type FeedbackDTO struct {
	PageNum  int     `json:"pageNum" binding:"required"`
	PageSize int     `json:"pageSize" binding:"required"`
	Keyword  *string `json:"keyword"`
}
