package dto

type UserResetPasswordDTO struct {
	Email            string `json:"email" binding:"required,email"`
	VerificationCode string `json:"verificationCode" binding:"required,verificationCode"`
	NewPassword      string `json:"newPassword" binding:"required,password"`
	RepeatPassword   string `json:"repeatPassword" binding:"required,password"`
}
