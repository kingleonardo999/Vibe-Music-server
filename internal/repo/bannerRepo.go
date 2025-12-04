package repo

import (
	"vibe-music-server/internal/model/entity"
	"vibe-music-server/internal/model/vo"
	"vibe-music-server/internal/pkg/db"
	"vibe-music-server/internal/pkg/result"
)

type BannerRepo struct{}

func NewBannerRepo() *BannerRepo {
	return &BannerRepo{}
}

func (b BannerRepo) GetPageBanners(data *result.PageResult[entity.Banner],
	status *uint8, index, size int) error {
	query := db.Get().Model(&entity.Banner{})
	// 可选条件
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	return query.
		Count(&data.Total).
		Offset(index).
		Limit(size).
		Order("id desc").
		Scan(&data.Items).Error
}

func (b BannerRepo) CreateBanner(banner *entity.Banner) error {
	return db.Get().Create(&banner).Error
}

func (b BannerRepo) GetBannerById(banner *entity.Banner, id uint64) error {
	return db.Get().First(banner, id).Error
}

func (b BannerRepo) GetBannerByIds(banner *[]entity.Banner, id []uint64) error {
	return db.Get().Where("id IN ?", id).First(banner).Error
}

func (b BannerRepo) UpdateBanner(banner *entity.Banner) error {
	return db.Get().Updates(banner).Error
}

func (b BannerRepo) DeleteBannerById(id uint64) error {
	return db.Get().Delete(&entity.Banner{}, id).Error
}

func (b BannerRepo) DeleteBannerByIds(id []uint64) error {
	return db.Get().Delete(&entity.Banner{}, id).Error
}

func (b BannerRepo) GetBannerList(banners *[]vo.BannerVO, limit int) error {
	return db.Get().Model(&entity.Banner{}).
		Select("id banner_id, banner_url").
		Where("status = ?", entity.BannerStatusEnable).
		Limit(limit).
		Order("id desc").
		Scan(&banners).Error
}
