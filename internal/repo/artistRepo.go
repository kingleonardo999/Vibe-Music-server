package repo

import (
	"vibe-music-server/internal/model/entity"
	"vibe-music-server/internal/model/vo"
	"vibe-music-server/internal/pkg/db"
	"vibe-music-server/internal/pkg/result"
)

type ArtistRepo struct{}

func NewArtistRepo() *ArtistRepo {
	return &ArtistRepo{}
}

func (a ArtistRepo) GetPageArtistsVO(data *result.PageResult[vo.ArtistVO], artistName *string,
	gender *uint8, area *string, startIndex, pageSize int) error {
	query := db.Get().Model(&entity.Artist{}).
		Select("id artist_id, name artist_name, avatar")

	// 可选条件
	if artistName != nil {
		query = query.Where("name LIKE ?", "%"+*artistName+"%")
	}
	if gender != nil {
		query = query.Where("gender = ?", *gender)
	}
	if area != nil {
		query = query.Where("area = ?", *area)
	}

	return query.
		Count(&data.Total).
		Offset(startIndex).
		Limit(pageSize).
		Scan(&data.Items).
		Error
}

func (a ArtistRepo) GetPageArtists(data *result.PageResult[entity.Artist], artistName *string,
	gender *uint8, area *string, startIndex, pageSize int) error {
	query := db.Get().Model(&entity.Artist{})
	// 可选条件
	if artistName != nil {
		query = query.Where("name LIKE ?", "%"+*artistName+"%")
	}
	if gender != nil {
		query = query.Where("gender = ?", *gender)
	}
	if area != nil {
		query = query.Where("area = ?", *area)
	}

	return query.
		Count(&data.Total).
		Offset(startIndex).
		Limit(pageSize).
		Order("id desc").
		Scan(&data.Items).
		Error
}

func (a ArtistRepo) GetAllArtistsName(data *[]vo.ArtistNameVO) error {
	return db.Get().Model(&entity.Artist{}).
		Select("id artist_id, name artist_name").
		Order("id desc").
		Scan(data).
		Error
}

func (a ArtistRepo) GetRandomArtists(data *[]vo.ArtistVO, limit int) error {
	return db.Get().Model(&entity.Artist{}).
		Select("id artist_id, name artist_name, avatar").
		Order("RAND()").
		Limit(limit).
		Scan(data).
		Error
}

func (a ArtistRepo) GetArtistDetail(data *vo.ArtistDetailVO, artistId uint64) error {
	if err := db.Get().Model(&entity.Artist{}).
		Select(`id artist_id, name artist_name, 
				gender, avatar, birth, area, introduction`).
		Where("id = ?", artistId).
		Scan(data).Error; err != nil {
		return err
	}
	// LikeStatus 默认0
	if err := db.Get().Model(&entity.Song{}).
		Select(`id song_id, name song_name,
			album, duration, cover_url, audio_url, release_time`).
		Where("artist_id = ?", artistId).
		Order("id desc").
		Scan(&data.Songs).Error; err != nil {
		return err
	}
	return nil
}

func (a ArtistRepo) GetArtistsCount(gender *uint8, area *string) (int64, error) {
	query := db.Get().Model(&entity.Artist{})
	// 可选条件
	if gender != nil {
		query = query.Where("gender = ?", *gender)
	}
	if area != nil {
		query = query.Where("area = ?", *area)
	}
	var count int64
	err := query.Count(&count).Error
	return count, err
}

func (a ArtistRepo) ExistArtistByName(name string) bool {
	query := db.Get().Model(&entity.Artist{}).Where("name = ?", name)
	var count int64
	query.Count(&count)
	return count > 0
}

func (a ArtistRepo) CreateArtist(artist *entity.Artist) error {
	return db.Get().Create(artist).Error
}

func (a ArtistRepo) SelectByName(artist *entity.Artist, name string) error {
	return db.Get().Model(&entity.Artist{}).Where("name = ?", name).First(artist).Error
}

func (a ArtistRepo) SelectById(artist *entity.Artist, artistId uint64) error {
	return db.Get().Model(&entity.Artist{}).Where("id = ?", artistId).First(artist).Error
}

func (a ArtistRepo) UpdateArtist(artist *entity.Artist, updateData any) error {
	return db.Get().Model(artist).Where("id = ?", artist.ID).Updates(updateData).Error
}

func (a ArtistRepo) DeleteArtistById(id uint64) error {
	return db.Get().Delete(&entity.Artist{}, id).Error
}

func (a ArtistRepo) GetAvatarsByIds(avatars *[]string, ids []uint64) error {
	return db.Get().Model(&entity.Artist{}).Select("avatar").Where("id IN ?", ids).Scan(avatars).Error
}

func (a ArtistRepo) DeleteArtistsByIds(ids []uint64) error {
	return db.Get().Where("id IN ?", ids).Delete(&entity.Artist{}).Error
}
