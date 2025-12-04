package entity

import "time"

type Artist struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	Name         string    `gorm:"size:100;not null;column:name"`
	Gender       uint8     `gorm:"type:tinyint;column:gender"` // 0-男 1-女
	Avatar       string    `gorm:"size:255;column:avatar"`
	Birth        time.Time `gorm:"type:date;column:birth"` // yyyy-MM-dd
	Area         string    `gorm:"size:100;column:area"`
	Introduction string    `gorm:"type:text;column:introduction"`
}

func (Artist) TableName() string { return "tb_artist" }
