package repo

import (
	"vibe-music-server/internal/model/entity"
	"vibe-music-server/internal/model/vo"
	"vibe-music-server/internal/pkg/db"
	"vibe-music-server/internal/pkg/result"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (u UserRepo) GetUserByName(user *entity.User, name string) error {
	return db.Get().Where("username = ?", name).First(user).Error
}

func (u UserRepo) GetUserByEmail(user *entity.User, email string) error {
	return db.Get().Where("email = ?", email).First(user).Error
}

func (u UserRepo) CreateUser(user *entity.User) error {
	return db.Get().Create(user).Error
}

func (u UserRepo) GetUserById(user *entity.User, id uint64) error {
	return db.Get().Where("id = ?", id).First(user).Error
}

func (u UserRepo) UpdateUser(user *entity.User) error {
	return db.Get().Model(user).Updates(user).Error
}

func (u UserRepo) DeleteUser(id uint64) error {
	return db.Get().Where("id = ?", id).Delete(&entity.User{}).Error
}

func (u UserRepo) DeleteUsers(ids []uint64) error {
	return db.Get().Where("id IN ?", ids).Delete(&entity.User{}).Error
}

func (u UserRepo) GetAllUsersCount(count *int64) error {
	return db.Get().Model(&entity.User{}).Count(count).Error
}

func (u UserRepo) GetUserByPhone(user *entity.User, phone string) error {
	return db.Get().Where("phone = ?", phone).First(user).Error
}

func (u UserRepo) GetAllUsers(data *result.PageResult[vo.UserManagementVO], username string, phone string, status uint8, index int, size int) error {
	query := db.Get().Model(&entity.User{}).
		Select("user_id, username, phone, email, status user_status, created_at, updated_at, user_avatar, introduction")
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}
	if status != 0 {
		query = query.Where("status = ?", status)
	}
	return query.Count(&data.Total).
		Order("created_at DESC").
		Offset(index).
		Limit(size).
		Scan(&data.Items).Error
}

func (u UserRepo) GetAvatarByIds(avatar *[]string, id []uint64) error {
	return db.Get().Model(&entity.User{}).Where("user_id IN ?", id).Pluck("user_avatar", avatar).Error
}
