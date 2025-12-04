package dto

type UserSearchDTO struct {
	PageNum  int    `json:"pageNum" binding:"required"`
	PageSize int    `json:"pageSize" binding:"required"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Status   uint8  `json:"status"` // 0-启用 1-禁用
}
