package vo

type PlaylistVO struct {
	PlaylistID uint64 `json:"playlistId"`
	Title      string `json:"title"`
	CoverURL   string `json:"coverUrl"`
}
