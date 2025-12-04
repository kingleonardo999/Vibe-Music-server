package dto

type PlaylistUpdateDTO struct {
	PlaylistID   uint64 `json:"playlistId"`
	Title        string `json:"title"`
	Introduction string `json:"introduction"`
	Style        string `json:"style"`
}
