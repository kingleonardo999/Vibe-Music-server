package dto

type CommentSongDTO struct {
	SongID  uint64 `json:"songId"`
	Content string `json:"content"`
}
