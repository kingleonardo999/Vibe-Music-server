package dto

type SongDTO struct {
	PageNum    int     `json:"pageNum" binding:"required"`
	PageSize   int     `json:"pageSize" binding:"required"`
	SongName   *string `json:"songName"`
	ArtistName *string `json:"artistName"`
	Album      *string `json:"album"`
}
