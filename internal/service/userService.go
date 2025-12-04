package service

import (
	"fmt"
	"time"
	"vibe-music-server/internal/model/dto"
	"vibe-music-server/internal/model/entity"
	"vibe-music-server/internal/model/vo"
	"vibe-music-server/internal/pkg/cache"
	"vibe-music-server/internal/pkg/result"
	"vibe-music-server/internal/pkg/result/consts"
	"vibe-music-server/internal/pkg/util"
	"vibe-music-server/internal/repo"
)

type UserService struct {
	userRepo     *repo.UserRepo
	emailService *EmailService
	minioService *MinioService
}

func NewUserService(userRepo *repo.UserRepo, emailService *EmailService, minioService *MinioService) *UserService {
	return &UserService{
		userRepo:     userRepo,
		emailService: emailService,
		minioService: minioService,
	}
}

func (u UserService) SendVerificationCode(email string) result.Result[result.Nil] {
	verificationCode := u.emailService.SendVerificationCodeEmail(email)
	if verificationCode == "" {
		return result.Error[result.Nil](consts.EmailSendFailed)
	}
	// 存入缓存，5分钟过期
	key := fmt.Sprintf("verificationCode:%s", email)
	expiration := 5 * time.Minute
	err := cache.SetWithExp(key, verificationCode, expiration)
	if err != nil {
		return result.Error[result.Nil](consts.InternalError)
	}
	return result.Success[result.Nil](consts.EmailSendSuccess)
}

func (u UserService) VerificationCode(email string, verificationCode string) bool {
	key := fmt.Sprintf("verificationCode:%s", email)
	code, err := cache.Get(key)
	if err != nil {
		return false
	}
	return code == verificationCode
}

func (u UserService) Register(userRegisterDTO *dto.UserRegisterDTO) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	// 验证验证码
	if !u.VerificationCode(userRegisterDTO.Email, userRegisterDTO.VerificationCode) {
		return retErr(consts.VerificationCode + consts.Error)
	}
	// 删除缓存中的验证码
	_ = cache.Del("verificationCode:" + userRegisterDTO.Email)
	// 判断用户是否存在
	var user entity.User
	if err := u.userRepo.GetUserByName(&user, userRegisterDTO.Username); err == nil {
		return retErr(consts.User + consts.AlreadyExists)
	}
	// 判断邮箱是否存在
	if err := u.userRepo.GetUserByEmail(&user, userRegisterDTO.Email); err == nil {
		return retErr(consts.Email + consts.AlreadyExists)
	}
	// 创建用户
	passwordHash, err := util.EncryptPassword(userRegisterDTO.Password)
	if err != nil {
		return retErr(consts.InternalError)
	}
	user = entity.User{
		Username:   userRegisterDTO.Username,
		Email:      userRegisterDTO.Email,
		Password:   passwordHash,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		Status:     entity.UserStatusEnable,
	}
	if err := u.userRepo.CreateUser(&user); err != nil {
		return retErr(consts.Register + consts.Failed)
	}
	// 清除用户相关缓存
	util.DeleteCacheByPattern("user:*")
	return retSuc(consts.Register + consts.Success)
}

func (u UserService) Login(userLoginDTO *dto.UserLoginDTO) result.Result[string] {
	retErr := result.Error[string]
	retSuc := result.SuccessWithData[string]
	var user entity.User
	if err := u.userRepo.GetUserByEmail(&user, userLoginDTO.Email); err != nil {
		return retErr(consts.InternalError)
	}
	if user.Email == "" {
		return retErr(consts.User + consts.NotExist)
	}
	if !util.ComparePassword(user.Password, userLoginDTO.Password) {
		return retErr(consts.Password + consts.Error)
	}
	if user.Status == entity.UserStatusDisable {
		return retErr(consts.User + consts.AccountLocked)
	}
	claims := util.Claims{
		Role:     consts.UserRole,
		UserId:   user.UserId,
		Username: user.Username,
	}
	token := util.GenerateToken(claims)

	// 将 token 存入缓存
	expiration := time.Hour * 12 // 12 小时后过期
	err := cache.SetWithExp(token, user.UserId, expiration)
	if err != nil {
		return retErr(consts.InternalError)
	}
	return retSuc(consts.Login+consts.Success, token)
}

// GetUserInfo claims 不能为nil
func (u UserService) GetUserInfo(claims *util.Claims) result.Result[vo.UserVO] {
	retErr := result.Error[vo.UserVO]
	retSuc := result.SuccessWithData[vo.UserVO]
	var user entity.User
	if err := u.userRepo.GetUserById(&user, claims.UserId); err != nil {
		return retErr(consts.InternalError)
	}
	if user.UserId == 0 {
		return retErr(consts.User + consts.NotExist)
	}
	userVO := vo.UserVO{
		UserID:       user.UserId,
		Username:     user.Username,
		Email:        user.Email,
		Phone:        *user.Phone,
		UserAvatar:   user.UserAvatar,
		Introduction: user.Introduction,
	}
	return retSuc(consts.Success, userVO)
}

func (u UserService) UpdateUserInfo(userDTO *dto.UserDTO) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	var user entity.User
	if err := u.userRepo.GetUserById(&user, userDTO.UserID); err != nil {
		return retErr(consts.User + consts.NotExist)
	}
	var userByName entity.User
	if err := u.userRepo.GetUserByName(&userByName, userDTO.Username); err == nil && userByName.UserId != user.UserId {
		return retErr(consts.Username + consts.AlreadyExists)
	}
	var userByEmail entity.User
	if err := u.userRepo.GetUserByEmail(&userByEmail, userDTO.Email); err == nil && userByEmail.UserId != user.UserId {
		return retErr(consts.Email + consts.AlreadyExists)
	}
	user.Username = userDTO.Username
	user.Email = userDTO.Email
	user.Phone = &userDTO.Phone
	user.Introduction = userDTO.Introduction
	user.UpdateTime = time.Now()
	if err := u.userRepo.UpdateUser(&user); err != nil {
		return retErr(consts.Update + consts.Failed)
	}
	util.DeleteCacheByPattern("user:*")
	return retSuc(consts.Update + consts.Success)
}

// UpdateUserAvatar claims 不能为nil
func (u UserService) UpdateUserAvatar(avatarUrl string, claims *util.Claims) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	var user entity.User
	if err := u.userRepo.GetUserById(&user, claims.UserId); err != nil {
		return retErr(consts.InternalError)
	}
	if user.UserId == 0 {
		return retErr(consts.User + consts.NotExist)
	}
	// 删除旧头像
	if user.UserAvatar != "" {
		err := u.minioService.DeleteFile(user.UserAvatar)
		if err != nil {
			return retErr(consts.Delete + consts.Failed)
		}
	}
	user.UserAvatar = avatarUrl
	user.UpdateTime = time.Now()
	if err := u.userRepo.UpdateUser(&user); err != nil {
		return retErr(consts.Update + consts.Failed)
	}
	util.DeleteCacheByPattern("user:*")
	return retSuc(consts.Update + consts.Success)
}

// UpdateUserPassword claims 不能为nil
func (u UserService) UpdateUserPassword(userPasswordDTO *dto.UserPasswordDTO, claims *util.Claims, token string) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	var user entity.User
	if err := u.userRepo.GetUserById(&user, claims.UserId); err != nil {
		return retErr(consts.InternalError)
	}
	if user.UserId == 0 {
		return retErr(consts.User + consts.NotExist)
	}
	if userPasswordDTO.OldPassword == userPasswordDTO.NewPassword {
		return retErr(consts.NewPasswordError)
	}
	if userPasswordDTO.OldPassword != userPasswordDTO.RepeatPassword {
		return retErr(consts.PasswordNotMatch)
	}
	if !util.ComparePassword(user.Password, userPasswordDTO.OldPassword) {
		return retErr(consts.OldPasswordError)
	}
	passwordHash, err := util.EncryptPassword(userPasswordDTO.NewPassword)
	if err != nil {
		return retErr(consts.InternalError)
	}
	user.Password = passwordHash
	user.UpdateTime = time.Now()
	if err := u.userRepo.UpdateUser(&user); err != nil {
		return retErr(consts.Update + consts.Failed)
	}
	// 删除缓存中的 token
	_ = cache.Del(token)
	util.DeleteCacheByPattern("user:*")
	return retSuc(consts.Update + consts.Success)
}

func (u UserService) ResetUserPassword(userResetPasswordDTO *dto.UserResetPasswordDTO) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	// 验证验证码
	if !u.VerificationCode(userResetPasswordDTO.Email, userResetPasswordDTO.VerificationCode) {
		return retErr(consts.VerificationCode + consts.Error)
	}
	// 删除缓存中的验证码
	_ = cache.Del("verificationCode:" + userResetPasswordDTO.Email)
	var user entity.User
	if err := u.userRepo.GetUserByEmail(&user, userResetPasswordDTO.Email); err != nil {
		return retErr(consts.InternalError)
	}
	if user.UserId == 0 {
		return retErr(consts.User + consts.NotExist)
	}
	// 校验两次密码
	if userResetPasswordDTO.NewPassword != userResetPasswordDTO.RepeatPassword {
		return retErr(consts.PasswordNotMatch)
	}
	passwordHash, err := util.EncryptPassword(userResetPasswordDTO.NewPassword)
	if err != nil {
		return retErr(consts.InternalError)
	}
	user.Password = passwordHash
	user.UpdateTime = time.Now()
	if err := u.userRepo.UpdateUser(&user); err != nil {
		return retErr(consts.Update + consts.Failed)
	}
	return retSuc(consts.Reset + consts.Success)
}

func (u UserService) Logout(token string) result.Result[result.Nil] {
	if err := cache.Del(token); err != nil {
		return result.Error[result.Nil](consts.InternalError)
	}
	util.DeleteCacheByPattern("user:*")
	util.DeleteCacheByPattern("playlist:*")
	util.DeleteCacheByPattern("song:*")
	util.DeleteCacheByPattern("favorite:*")
	util.DeleteCacheByPattern("artist:*")
	return result.Success[result.Nil](consts.Logout + consts.Success)
}

// DeleteAccount claims 不能为nil
func (u UserService) DeleteAccount(claims *util.Claims, token string) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	var user entity.User
	if err := u.userRepo.GetUserById(&user, claims.UserId); err != nil {
		return retErr(consts.InternalError)
	}
	if user.UserId == 0 {
		return retErr(consts.User + consts.NotExist)
	}
	// 删除用户头像
	if user.UserAvatar != "" {
		err := u.minioService.DeleteFile(user.UserAvatar)
		if err != nil {
			return retErr(consts.Delete + consts.Failed)
		}
	}
	if err := u.userRepo.DeleteUser(claims.UserId); err != nil {
		return retErr(consts.Delete + consts.Failed)
	}
	// 删除缓存中的 token
	_ = cache.Del(token)
	util.DeleteCacheByPattern("user:*")
	return retSuc(consts.Delete + consts.Success)
}

func (u UserService) GetAllUsersCount() result.Result[int64] {
	var count int64
	_ = u.userRepo.GetAllUsersCount(&count)
	return result.SuccessWithData[int64](consts.Success, count)
}

func (u UserService) GetAllUsers(userSearchDTO *dto.UserSearchDTO) result.Result[result.PageResult[vo.UserManagementVO]] {
	retErr := result.Error[result.PageResult[vo.UserManagementVO]]
	retSuc := result.SuccessWithData[result.PageResult[vo.UserManagementVO]]
	pageNum := userSearchDTO.PageNum
	pageSize := userSearchDTO.PageSize
	startIndex := (pageNum - 1) * pageSize
	var data result.PageResult[vo.UserManagementVO]
	templateKey := fmt.Sprintf("user:getAllUsers:%v-%v-%v-%v-%v", userSearchDTO.Username, userSearchDTO.Phone, userSearchDTO.Status, startIndex, pageSize)
	if !util.GetCache(templateKey, &data) {
		if err := u.userRepo.GetAllUsers(&data, userSearchDTO.Username, userSearchDTO.Phone, userSearchDTO.Status, startIndex, pageSize); err != nil {
			return retErr(consts.InternalError)
		}
		util.SetCache(templateKey, data)
	}
	if len(data.Items) == 0 {
		return retErr(consts.DataNotFound)
	}
	return retSuc(consts.Success, data)
}

func (u UserService) AddUser(userAddDTO *dto.UserAddDTO) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	// 判断用户是否存在
	var userByName entity.User
	if err := u.userRepo.GetUserByName(&userByName, userAddDTO.Username); err == nil {
		return retErr(consts.Username + consts.AlreadyExists)
	}
	var userByEmail entity.User
	if err := u.userRepo.GetUserByEmail(&userByEmail, userAddDTO.Email); err == nil {
		return retErr(consts.Email + consts.AlreadyExists)
	}
	var userByPhone entity.User
	if err := u.userRepo.GetUserByPhone(&userByPhone, userAddDTO.Phone); err == nil {
		return retErr(consts.Phone + consts.AlreadyExists)
	}
	// 创建用户
	var user entity.User
	passwordHash, err := util.EncryptPassword(userAddDTO.Password)
	if err != nil {
		return retErr(consts.InternalError)
	}
	user = entity.User{
		Username:     userAddDTO.Username,
		Email:        userAddDTO.Email,
		Phone:        &userAddDTO.Phone,
		Password:     passwordHash,
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
		Status:       entity.UserStatus(userAddDTO.Status),
		Introduction: userAddDTO.Introduction,
	}
	if err := u.userRepo.CreateUser(&user); err != nil {
		return retErr(consts.Add + consts.Failed)
	}
	// 清除用户相关缓存
	util.DeleteCacheByPattern("user:*")
	return retSuc(consts.Add + consts.Success)
}

func (u UserService) UpdateUser(userDTO *dto.UserDTO) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	var user entity.User
	if err := u.userRepo.GetUserById(&user, userDTO.UserID); err != nil {
		return retErr(consts.User + consts.NotExist)
	}
	var userByName entity.User
	if err := u.userRepo.GetUserByName(&userByName, userDTO.Username); err == nil && userByName.UserId != user.UserId {
		return retErr(consts.Username + consts.AlreadyExists)
	}
	var userByEmail entity.User
	if err := u.userRepo.GetUserByEmail(&userByEmail, userDTO.Email); err == nil && userByEmail.UserId != user.UserId {
		return retErr(consts.Email + consts.AlreadyExists)
	}
	var userByPhone entity.User
	if err := u.userRepo.GetUserByPhone(&userByPhone, userDTO.Phone); err == nil && userByPhone.UserId != user.UserId {
		return retErr(consts.Phone + consts.AlreadyExists)
	}
	user.Username = userDTO.Username
	user.Email = userDTO.Email
	user.Phone = &userDTO.Phone
	user.Introduction = userDTO.Introduction
	user.UpdateTime = time.Now()
	if err := u.userRepo.UpdateUser(&user); err != nil {
		return retErr(consts.Update + consts.Failed)
	}
	util.DeleteCacheByPattern("user:*")
	return retSuc(consts.Update + consts.Success)
}

func (u UserService) UpdateUserStatus(userId uint64, userStatus entity.UserStatus) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	var user entity.User
	if err := u.userRepo.GetUserById(&user, userId); err != nil {
		return retErr(consts.User + consts.NotExist)
	}
	user.Status = userStatus
	user.UpdateTime = time.Now()
	if err := u.userRepo.UpdateUser(&user); err != nil {
		return retErr(consts.Update + consts.Failed)
	}
	util.DeleteCacheByPattern("user:*")
	return retSuc(consts.Update + consts.Success)
}

func (u UserService) DeleteUser(userId uint64) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	var user entity.User
	if err := u.userRepo.GetUserById(&user, userId); err != nil {
		return retErr(consts.InternalError)
	}
	if user.UserId == 0 {
		return retErr(consts.User + consts.NotExist)
	}
	// 删除用户头像
	if user.UserAvatar != "" {
		err := u.minioService.DeleteFile(user.UserAvatar)
		if err != nil {
			return retErr(consts.Delete + consts.Failed)
		}
	}
	if err := u.userRepo.DeleteUser(userId); err != nil {
		return retErr(consts.Delete + consts.Failed)
	}
	util.DeleteCacheByPattern("user:*")
	return retSuc(consts.Delete + consts.Success)
}

func (u UserService) DeleteUsers(userIds []uint64) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	var avatarUrls []string
	if err := u.userRepo.GetAvatarByIds(&avatarUrls, userIds); err != nil {
		return retErr(consts.InternalError)
	}
	// 删除用户头像
	for _, avatarUrl := range avatarUrls {
		if avatarUrl != "" {
			err := u.minioService.DeleteFile(avatarUrl)
			if err != nil {
				return retErr(consts.Delete + consts.Failed)
			}
		}
	}
	if err := u.userRepo.DeleteUsers(userIds); err != nil {
		return retErr(consts.Delete + consts.Failed)
	}
	util.DeleteCacheByPattern("user:*")
	return retSuc(consts.Delete + consts.Success)
}
