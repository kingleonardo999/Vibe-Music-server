package vo

type PlaylistDetailVO struct {
	PlaylistID   uint64      `json:"playlistId"`
	Title        string      `json:"title"`
	CoverURL     string      `json:"coverUrl"`
	Introduction string      `json:"introduction"`
	Songs        []SongVO    `json:"songs"`      // 歌曲简要列表
	LikeStatus   uint8       `json:"likeStatus"` // 0-默认 1-喜欢
	Comments     []CommentVO `json:"comments"`   // 评论列表
}
