package dto

import "time"

type ArtistUpdateDTO struct {
	ArtistID     uint64    `json:"artistId"`
	ArtistName   string    `json:"artistName"`
	Gender       uint8     `json:"gender"`                         // 0-男 1-女
	Birth        time.Time `json:"birth" time_format:"2006-01-02"` // yyyy-MM-dd
	Area         string    `json:"area"`
	Introduction string    `json:"introduction"`
}
