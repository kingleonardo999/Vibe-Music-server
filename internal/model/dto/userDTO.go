package dto

type UserDTO struct {
	UserID       uint64 `json:"userId"`
	Username     string `json:"username" binding:"required,username"`
	Phone        string `json:"phone" binding:"phone"`
	Email        string `json:"email" binding:"required,email"`
	Introduction string `json:"introduction" binding:"max=100"`
}
