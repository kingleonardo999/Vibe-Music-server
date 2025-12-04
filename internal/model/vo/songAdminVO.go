package vo

import "time"

type SongAdminVO struct {
	SongID      uint64    `json:"songId"`
	ArtistName  string    `json:"artistName"`
	SongName    string    `json:"songName"`
	Album       string    `json:"album"`
	Lyric       string    `json:"lyric"`
	Duration    string    `json:"duration"`
	Style       string    `json:"style"`
	CoverURL    string    `json:"coverUrl"`
	AudioURL    string    `json:"audioUrl"`
	ReleaseTime time.Time `json:"releaseTime" time_format:"2006-01-02"`
}
