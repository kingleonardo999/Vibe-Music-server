package repo

import (
	"vibe-music-server/internal/model/entity"
	"vibe-music-server/internal/pkg/db"
)

type FavoriteRepo struct{}

func NewFavoriteRepo() *FavoriteRepo {
	return &FavoriteRepo{}
}

func (f FavoriteRepo) GetFavoriteSongIds(data *[]uint64, userId uint64) error {
	query := db.Get().Model(&entity.Favorite{}).
		Select("song_id").
		Where("user_id = ? AND type = ?", userId, entity.FavoriteTypeSong).
		Order("song_id").
		Scan(data)
	return query.Error
}

func (f FavoriteRepo) GetFavoritePlaylistIds(data *[]uint64, userId uint64) error {
	query := db.Get().Model(&entity.Favorite{}).
		Select("playlist_id").
		Where("user_id = ? AND type = ?", userId, entity.FavoriteTypePlaylist).
		Order("playlist_id").
		Scan(data)
	return query.Error
}

func (f FavoriteRepo) GetFavoritePlaylists(data *[]entity.Favorite, userId uint64) error {
	query := db.Get().Model(&entity.Favorite{}).
		Where("user_id = ? AND type = ?", userId, entity.FavoriteTypePlaylist).
		Order("playlist_id").
		Scan(data)
	return query.Error
}

func (f FavoriteRepo) GetFavoriteSongs(data *[]entity.Song, userId uint64) error {
	query := db.Get().Model(&entity.Favorite{}).
		Where("user_id = ? AND type = ?", userId, entity.FavoriteTypeSong).
		Order("song_id").
		Scan(data)
	return query.Error
}

func (f FavoriteRepo) IsFavoritePlaylist(isFavorite *uint8, userId uint64, playlistId uint64) error {
	query := db.Get().Model(&entity.Favorite{}).
		Select("COUNT(1)").
		Where("user_id = ? AND playlist_id = ? AND type = ?", userId, playlistId, entity.FavoriteTypePlaylist).
		Scan(isFavorite)
	return query.Error
}

func (f FavoriteRepo) IsFavoriteSong(isFavorite *uint8, userId uint64, songId uint64) error {
	query := db.Get().Model(&entity.Favorite{}).
		Select("COUNT(1)").
		Where("user_id = ? AND song_id = ? AND type = ?", userId, songId, entity.FavoriteTypeSong).
		Scan(isFavorite)
	return query.Error
}

func (f FavoriteRepo) AddFavorite(favorite *entity.Favorite) error {
	query := db.Get().Create(favorite)
	return query.Error
}

func (f FavoriteRepo) DeleteFavoritePlaylist(userId uint64, playlistId uint64) error {
	query := db.Get().Where("user_id = ? AND playlist_id = ? AND type = ?", userId, playlistId, entity.FavoriteTypePlaylist).
		Delete(&entity.Favorite{})
	return query.Error
}

func (f FavoriteRepo) DeleteFavoriteSong(userId uint64, songId uint64) error {
	query := db.Get().Where("user_id = ? AND song_id = ? AND type = ?", userId, songId, entity.FavoriteTypeSong).
		Delete(&entity.Favorite{})
	return query.Error
}
