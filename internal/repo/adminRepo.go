package repo

import (
	"vibe-music-server/internal/model/entity"
	"vibe-music-server/internal/pkg/db"
)

type AdminRepo struct{}

func NewAdminRepo() *AdminRepo {
	return &AdminRepo{}
}

func (a AdminRepo) SelectByUsername(username string) (*entity.Admin, error) {
	var admin entity.Admin
	err := db.Get().Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (a AdminRepo) Insert(admin *entity.Admin) error {
	return db.Get().Create(admin).Error
}
