package dto

type BannerDTO struct {
	PageNum  int    `json:"pageNum" binding:"required"`
	PageSize int    `json:"pageSize" binding:"required"`
	Status   *uint8 `json:"status"` // 0-启用 1-禁用
}
