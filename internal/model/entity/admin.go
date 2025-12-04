package entity

type Admin struct {
	AdminId  uint64 `gorm:"primaryKey;autoIncrement;column:id"`
	Username string `gorm:"unique;not null;column:username"`
	Password string `gorm:"not null;column:password"` // 存哈希后长度
}

func (Admin) TableName() string {
	return "tb_admin"
}
