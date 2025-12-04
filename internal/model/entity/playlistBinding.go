package entity

type PlaylistBinding struct {
	PlaylistID uint64 `gorm:"primaryKey;column:playlist_id"`
	SongID     uint64 `gorm:"primaryKey;column:song_id"`
}

// 联合主键：歌单+歌曲

func (PlaylistBinding) TableName() string { return "tb_playlist_binding" }
