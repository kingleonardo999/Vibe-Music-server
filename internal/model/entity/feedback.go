package entity

import "time"

type Feedback struct {
	ID         uint      `gorm:"primaryKey;autoIncrement;column:id"`
	UserID     uint64    `gorm:"index;not null;column:user_id"`
	Feedback   string    `gorm:"type:text;not null;column:feedback"`
	CreateTime time.Time `gorm:"type:datetime;not null;column:create_time"`
}

func (Feedback) TableName() string { return "tb_feedback" }
