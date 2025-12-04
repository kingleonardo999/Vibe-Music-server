package service

import (
	"errors"
	"gorm.io/gorm"
	"time"
	"vibe-music-server/internal/model/dto"
	"vibe-music-server/internal/model/entity"
	"vibe-music-server/internal/pkg/cache"
	"vibe-music-server/internal/pkg/result"
	"vibe-music-server/internal/pkg/result/consts"
	"vibe-music-server/internal/pkg/util"
	"vibe-music-server/internal/repo"
)

type AdminService struct {
	adminRepo *repo.AdminRepo
}

func NewAdminService(adminRepo *repo.AdminRepo) *AdminService {
	return &AdminService{
		adminRepo: adminRepo,
	}
}

// Register 管理员注册
func (a AdminService) Register(adminDTO *dto.AdminDTO) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	admin, err := a.adminRepo.SelectByUsername(adminDTO.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return retErr(consts.InternalError)
	}
	if admin != nil {
		return retErr(consts.Username + consts.AlreadyExists)
	}
	encryptedPassword, err := util.EncryptPassword(adminDTO.Password)
	if err != nil {
		return retErr(consts.InternalError)
	}
	admin = &entity.Admin{
		Username: adminDTO.Username,
		Password: encryptedPassword,
	}
	err = a.adminRepo.Insert(admin)
	if err != nil {
		return retErr(consts.Register + consts.Failed)
	}
	return retSuc(consts.Register + consts.Success)
}

// Login 管理员登录
func (a AdminService) Login(adminDTO *dto.AdminDTO) result.Result[string] {
	retErr := result.Error[string]
	retSuc := result.SuccessWithData[string]
	admin, err := a.adminRepo.SelectByUsername(adminDTO.Username)
	if err != nil {
		return retErr(consts.InternalError)
	}
	if admin == nil {
		return retErr(consts.User + consts.NotExist)
	}
	if !util.ComparePassword(admin.Password, adminDTO.Password) {
		return retErr(consts.User + consts.Invalid)
	}
	claims := util.Claims{
		Role:     consts.AdminRole,
		UserId:   admin.AdminId,
		Username: admin.Username,
	}
	token := util.GenerateToken(claims)

	// 将 token 存入缓存
	expiration := time.Hour * 12 // 12 小时后过期
	err = cache.SetWithExp(token, admin.AdminId, expiration)
	if err != nil {
		return retErr(consts.InternalError)
	}
	return retSuc(consts.Login+consts.Success, token)
}

// Logout 管理员登出
func (a AdminService) Logout(token string) result.Result[result.Nil] {
	ret := cache.Del(token)
	if ret != nil {
		return result.Error[result.Nil](consts.InternalError)
	}
	return result.Success[result.Nil](consts.Logout + consts.Success)
}
