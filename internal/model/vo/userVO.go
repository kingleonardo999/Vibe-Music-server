package vo

type UserVO struct {
	UserID       uint64 `json:"userId"`
	Username     string `json:"username"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	UserAvatar   string `json:"userAvatar"`
	Introduction string `json:"introduction"`
}
