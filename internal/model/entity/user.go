package entity

import "time"

type UserStatus uint8

const (
	UserStatusEnable  UserStatus = 0 // 启用
	UserStatusDisable UserStatus = 1 // 禁用
)

type User struct {
	UserId       uint64     `gorm:"primaryKey;autoIncrement;column:id"`
	Username     string     `gorm:"size:16;unique;not null;column:username"`
	Password     string     `gorm:"size:128;not null;column:password"` // 存哈希
	Phone        *string    `gorm:"size:11;column:phone"`              // 可为空
	Email        string     `gorm:"size:100;not null;column:email"`
	UserAvatar   string     `gorm:"size:500;column:user_avatar"`
	Introduction string     `gorm:"size:100;column:introduction"`
	CreateTime   time.Time  `gorm:"type:datetime;not null;column:create_time"`
	UpdateTime   time.Time  `gorm:"type:datetime;not null;column:update_time"`
	Status       UserStatus `gorm:"type:tinyint;not null;column:status"` // 0-启用 1-禁用
}

func (User) TableName() string { return "tb_user" }
