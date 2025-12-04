package vo

type ArtistVO struct {
	ArtistID   uint64 `json:"artistId"`
	ArtistName string `json:"artistName"`
	Avatar     string `json:"avatar"`
}
