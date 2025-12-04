package dto

import "time"

type SongUpdateDTO struct {
	SongID      uint64    `json:"songId"`
	ArtistID    uint64    `json:"artistId"`
	SongName    string    `json:"songName"`
	Album       string    `json:"album"`
	Style       string    `json:"style"`
	ReleaseTime time.Time `json:"releaseTime" time_format:"2006-01-02"` // yyyy-MM-dd
}
