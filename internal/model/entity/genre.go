package entity

type Genre struct {
	SongID  uint64 `gorm:"primaryKey;column:song_id"`
	StyleID uint64 `gorm:"column:style_id"`
}

func (Genre) TableName() string { return "tb_genre" }
