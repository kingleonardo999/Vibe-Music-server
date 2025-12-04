package repo

import (
	"vibe-music-server/internal/model/entity"
	"vibe-music-server/internal/pkg/db"
)

type StyleRepo struct{}

func NewStyleRepo() *StyleRepo {
	return &StyleRepo{}
}

func (s StyleRepo) GetStyleIdsByNames(ids *[]uint64, names []string) error {
	query := db.Get().Model(&entity.Style{}).
		Select("id").
		Where("name IN ?", names).
		Pluck("id", ids)
	return query.Error
}
