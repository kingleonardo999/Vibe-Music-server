package dto

type UserRegisterDTO struct {
	Username         string `json:"username" binding:"required,username"`
	Password         string `json:"password" binding:"required,password"`
	Email            string `json:"email" binding:"required,email"`
	VerificationCode string `json:"verificationCode" binding:"required,verificationCode"`
}
