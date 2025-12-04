package dto

type PlaylistDTO struct {
	PageNum  int     `json:"pageNum" binding:"required"`
	PageSize int     `json:"pageSize" binding:"required"`
	Title    *string `json:"title"`
	Style    *string `json:"style"`
}
