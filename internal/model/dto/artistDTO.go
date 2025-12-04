package dto

type ArtistDTO struct {
	PageNum    int     `json:"pageNum" binding:"required"`
	PageSize   int     `json:"pageSize" binding:"required"`
	ArtistName *string `json:"artistName"`
	Gender     *uint8  `json:"gender"` // 0-男 1-女 2-组合
	Area       *string `json:"area"`
}
