package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"vibe-music-server/internal/model/dto"
	"vibe-music-server/internal/pkg/result"
	"vibe-music-server/internal/pkg/result/consts"
	"vibe-music-server/internal/pkg/util"
	"vibe-music-server/internal/pkg/validate"
	"vibe-music-server/internal/service"
)

type UserCtrl struct {
	userService  *service.UserService
	minioService *service.MinioService
}

func NewUserCtrl(userService *service.UserService, minioService *service.MinioService) *UserCtrl {
	return &UserCtrl{
		userService:  userService,
		minioService: minioService,
	}
}

func (u *UserCtrl) SendVerificationCode(c *gin.Context) {
	email := c.Request.FormValue("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	if err := validate.Validate.Var(email, "email"); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, u.userService.SendVerificationCode(email))
}

func (u *UserCtrl) Register(c *gin.Context) {
	var userRegisterDTO dto.UserRegisterDTO
	if err := c.ShouldBindJSON(&userRegisterDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	if !u.userService.VerificationCode(userRegisterDTO.Email, userRegisterDTO.VerificationCode) {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.VerificationCode+consts.Invalid))
	}
	c.JSON(http.StatusOK, u.userService.Register(&userRegisterDTO))
}

func (u *UserCtrl) Login(c *gin.Context) {
	var userLoginDTO dto.UserLoginDTO
	if err := c.ShouldBindJSON(&userLoginDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[string](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, u.userService.Login(&userLoginDTO))
}

// GetUserInfo
// need authMiddleware
func (u *UserCtrl) GetUserInfo(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, result.Error[result.Nil](consts.NotLogin))
		return
	}
	c.JSON(http.StatusOK, u.userService.GetUserInfo(claims.(*util.Claims)))
}

// UpdateUserInfo
// need authMiddleware
func (u *UserCtrl) UpdateUserInfo(c *gin.Context) {
	var userDTO dto.UserDTO
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	claimsI, _ := c.Get("claims")
	claims := claimsI.(*util.Claims)
	if claims.UserId != userDTO.UserID {
		c.JSON(http.StatusForbidden, result.Error[result.Nil](consts.NoPermission))
		return
	}
	c.JSON(http.StatusOK, u.userService.UpdateUserInfo(&userDTO))
}

// UpdateUserAvatar
// need authMiddleware
func (u *UserCtrl) UpdateUserAvatar(c *gin.Context) {
	avatar, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[string](consts.InvalidParams))
		return
	}
	avatarUrl, err := u.minioService.UploadFile(avatar, "users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.Error[string](consts.FileUpload+consts.Failed))
		return
	}
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, result.Error[result.Nil](consts.NotLogin))
		return
	}
	c.JSON(http.StatusOK, u.userService.UpdateUserAvatar(avatarUrl, claims.(*util.Claims)))
}

// UpdateUserPassword
// need authMiddleware
func (u *UserCtrl) UpdateUserPassword(c *gin.Context) {
	var userPasswordDTO dto.UserPasswordDTO
	if err := c.ShouldBindJSON(&userPasswordDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, result.Error[result.Nil](consts.NotLogin))
		return
	}
	token := c.GetHeader("Authorization")
	c.JSON(http.StatusOK, u.userService.UpdateUserPassword(&userPasswordDTO, claims.(*util.Claims), token))
}

func (u *UserCtrl) ResetUserPassword(c *gin.Context) {
	var userResetPasswordDTO dto.UserResetPasswordDTO
	if err := c.ShouldBindJSON(&userResetPasswordDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	if !u.userService.VerificationCode(userResetPasswordDTO.Email, userResetPasswordDTO.VerificationCode) {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.VerificationCode+consts.Invalid))
		return
	}
	c.JSON(http.StatusOK, u.userService.ResetUserPassword(&userResetPasswordDTO))
}

// Logout
// need authMiddleware
func (u *UserCtrl) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, u.userService.Logout(token))
}

// DeleteAccount
// need authMiddleware
func (u *UserCtrl) DeleteAccount(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, result.Error[result.Nil](consts.NotLogin))
		return
	}
	token := c.GetHeader("Authorization")
	c.JSON(http.StatusOK, u.userService.DeleteAccount(claims.(*util.Claims), token))
}
