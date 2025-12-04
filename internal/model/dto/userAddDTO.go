package dto

type UserAddDTO struct {
	Username     string `json:"username" binding:"required,username"`
	Password     string `json:"password" binding:"required,password"`
	Phone        string `json:"phone" binding:"phone"`
	Email        string `json:"email" binding:"required,email"`
	Introduction string `json:"introduction"`
	Status       uint8  `json:"status"` // 0-启用 1-禁用
}
