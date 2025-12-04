package entity

import "time"

type FavoriteType uint8

const (
	FavoriteTypeSong     FavoriteType = 0 // 歌曲
	FavoriteTypePlaylist FavoriteType = 1 // 歌单
)

type Favorite struct {
	ID         uint64       `gorm:"primaryKey;autoIncrement;column:id"`
	UserID     uint64       `gorm:"index;not null;column:user_id"`
	Type       FavoriteType `gorm:"type:tinyint;not null;column:type"` // 0-歌曲 1-歌单
	SongID     *uint64      `gorm:"index;column:song_id"`              // 收藏歌曲时非空
	PlaylistID *uint64      `gorm:"index;column:playlist_id"`          // 收藏歌单时非空
	CreateTime time.Time    `gorm:"type:datetime;not null;column:create_time"`
}

func (Favorite) TableName() string { return "tb_user_favorite" }
