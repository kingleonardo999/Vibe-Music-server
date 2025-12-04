package dto

type SongAndArtistDTO struct {
	PageNum  int     `json:"pageNum" binding:"required"`
	PageSize int     `json:"pageSize" binding:"required"`
	ArtistID *uint64 `json:"artistId"`
	SongName *string `json:"songName"`
	Album    *string `json:"album"`
}
