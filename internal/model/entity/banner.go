package entity

type BannerStatus uint8

const (
	BannerStatusEnable  BannerStatus = 0 // 启用
	BannerStatusDisable BannerStatus = 1 // 禁用
)

type Banner struct {
	ID        uint64       `gorm:"primaryKey;autoIncrement;column:id"`
	BannerURL string       `gorm:"size:255;not null;column:banner_url"`
	Status    BannerStatus `gorm:"type:tinyint;not null;column:status"`
}

func (Banner) TableName() string { return "tb_banner" }
