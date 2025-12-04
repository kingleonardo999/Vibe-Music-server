package vo

import "time"

type UserManagementVO struct {
	UserID       uint64    `json:"userId"`
	Username     string    `json:"username"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	UserAvatar   string    `json:"userAvatar"`
	Introduction string    `json:"introduction"`
	CreateTime   time.Time `json:"createTime"`
	UpdateTime   time.Time `json:"updateTime"`
	UserStatus   uint8     `json:"userStatus"` // 0-启用 1-禁用
}
