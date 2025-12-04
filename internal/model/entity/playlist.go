package entity

type Playlist struct {
	ID           uint   `gorm:"primaryKey;autoIncrement;column:id"`
	Title        string `gorm:"size:200;not null;column:title"`
	CoverURL     string `gorm:"size:500;column:cover_url"`
	Introduction string `gorm:"type:text;column:introduction"`
	Style        string `gorm:"size:100;column:style"`
}

func (Playlist) TableName() string { return "tb_playlist" }
