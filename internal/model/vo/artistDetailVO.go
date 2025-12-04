package vo

import (
	"time"
)

type ArtistDetailVO struct {
	ArtistID     uint64    `json:"artistId"`
	ArtistName   string    `json:"artistName"`
	Gender       uint8     `json:"gender"` // 0-男 1-女
	Avatar       string    `json:"avatar"`
	Birth        time.Time `json:"birth"    time_format:"2006-01-02"` // 仅日期
	Area         string    `json:"area"`
	Introduction string    `json:"introduction"`
	Songs        []SongVO  `json:"songs"` // 内嵌歌曲简要信息
}
