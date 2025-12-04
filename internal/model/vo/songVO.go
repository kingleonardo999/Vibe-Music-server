package vo

import "time"

type SongVO struct {
	SongID      uint64    `json:"songId"`
	SongName    string    `json:"songName"`
	ArtistName  string    `json:"artistName"`
	Album       string    `json:"album"`
	Duration    string    `json:"duration"`
	CoverURL    string    `json:"coverUrl"`
	AudioURL    string    `json:"audioUrl"`
	LikeStatus  uint8     `json:"likeStatus"` // 0-默认 1-喜欢
	ReleaseTime time.Time `json:"releaseTime" time_format:"2006-01-02"`
}
