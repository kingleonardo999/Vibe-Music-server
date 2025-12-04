package entity

import "time"

type CommentType uint8

const (
	CommentTypeSong     CommentType = 0
	CommentTypePlaylist CommentType = 1
)

type Comment struct {
	ID         uint        `gorm:"primaryKey;autoIncrement;column:id"`
	UserID     uint64      `gorm:"index;not null;column:user_id"`
	SongID     *uint64     `gorm:"index;column:song_id"`     // 歌曲评论时非空
	PlaylistID *uint64     `gorm:"index;column:playlist_id"` // 歌单评论时非空
	Content    string      `gorm:"type:text;not null;column:content"`
	CreateTime time.Time   `gorm:"type:datetime;not null;column:create_time"`
	Type       CommentType `gorm:"type:tinyint;not null;column:type"` // 0-歌曲 1-歌单
	LikeCount  uint        `gorm:"default:0;column:like_count"`
}

func (Comment) TableName() string { return "tb_comment" }
