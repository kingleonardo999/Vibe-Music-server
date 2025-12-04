package entity

type Style struct {
	ID   uint   `gorm:"primaryKey;autoIncrement;column:id"`
	Name string `gorm:"size:100;not null;column:name"`
}

func (Style) TableName() string { return "tb_style" }
