package dto

type UserPasswordDTO struct {
	OldPassword    string `json:"oldPassword" binding:"required,password"`
	NewPassword    string `json:"newPassword" binding:"required,password"`
	RepeatPassword string `json:"repeatPassword" binding:"required,password"`
}
