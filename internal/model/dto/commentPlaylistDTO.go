package dto

type CommentPlaylistDTO struct {
	PlaylistID uint64 `json:"playlistId"`
	Content    string `json:"content"`
}
