package repo

import (
	"vibe-music-server/internal/model/entity"
	"vibe-music-server/internal/pkg/db"
)

type GenreRepo struct{}

func NewGenreRepo() *GenreRepo {
	return &GenreRepo{}
}

func (r GenreRepo) CreateGenre(genre *entity.Genre) error {
	return db.Get().Create(&genre).Error
}

func (r GenreRepo) DeleteGenresBySongId(id uint64) error {
	return db.Get().Where("song_id = ?", id).Delete(&entity.Genre{}).Error
}

func (r GenreRepo) DeleteGenresBySongIds(ids []uint64) error {
	return db.Get().Where("song_id IN ?", ids).Delete(&entity.Genre{}).Error
}
