package repo

import (
	"gorm.io/gorm"
	"vibe-music-server/internal/model/dto"
	"vibe-music-server/internal/model/entity"
	"vibe-music-server/internal/model/vo"
	"vibe-music-server/internal/pkg/db"
	"vibe-music-server/internal/pkg/result"
)

type PlaylistRepo struct{}

func NewPlaylistRepo() *PlaylistRepo {
	return &PlaylistRepo{}
}

func (p PlaylistRepo) GetAllPlaylists(data *result.PageResult[vo.PlaylistVO], title, style *string, index, size int) error {
	query := db.Get().Model(&entity.Playlist{}).
		Select("id playlist_id, title, cover_url")
	if title != nil {
		query = query.Where("title LIKE ?", "%"+*title+"%")
	}
	if style != nil {
		query = query.Where("style = ?", *style)
	}
	return query.Count(&data.Total).
		Offset(index).
		Limit(size).
		Scan(&data.Items).Error
}

func (p PlaylistRepo) GetAllPlaylistsInfo(data *result.PageResult[entity.Playlist], title, style *string, index, size int) error {
	query := db.Get().Model(&entity.Playlist{})
	if title != nil {
		query = query.Where("title LIKE ?", "%"+*title+"%")
	}
	if style != nil {
		query = query.Where("style = ?", *style)
	}
	return query.Count(&data.Total).
		Offset(index).
		Limit(size).
		Find(&data.Items).Error
}

func (p PlaylistRepo) GetRandomPlaylists(data *[]vo.PlaylistVO, limit int) error {
	query := db.Get().Model(&entity.Playlist{}).
		Select("id playlist_id, title, cover_url").Order("RAND()").
		Limit(limit)
	return query.Scan(data).Error
}

func (p PlaylistRepo) GetStylesByIds(styles *[]string, ids []uint64) error {
	query := db.Get().Model(&entity.Playlist{}).
		Distinct("style").
		Where("id IN ?", ids).
		Order("style").
		Pluck("style", styles)
	return query.Error
}

func (p PlaylistRepo) GetRecommendedPlaylistsByStyles(data *[]vo.PlaylistVO, styles []string, ids []uint64, limit int) error {
	query := db.Get().Model(&entity.Playlist{}).
		Select("id playlist_id, title, cover_url").
		Where("style IN ?", styles).
		Where("id Not IN ?", ids).
		Order("RANDOM()").
		Limit(limit)
	return query.Scan(data).Error
}

func (p PlaylistRepo) GetPlaylistDetail(data *vo.PlaylistDetailVO, id uint64) error {
	query := db.Get().Model(&entity.Playlist{}).
		Select("id playlist_id, title, cover_url, introduction").
		Where("id = ?", id).
		Scan(data)
	songQuery := db.Get().Model(&entity.PlaylistBinding{}).
		Where("playlist_id = ?", id).Scan(&data.Songs)
	commentQuery := db.Get().Model(&entity.Comment{}).
		Where("type = ? AND playlist_id = ?", entity.CommentTypePlaylist, id).Scan(&data.Comments)
	switch {
	case query.Error != nil:
		return query.Error
	case songQuery.Error != nil:
		return songQuery.Error
	case commentQuery.Error != nil:
		return commentQuery.Error
	}
	return nil
}

func (p PlaylistRepo) GetAllPlaylistsCount(count *int64, style *string) error {
	query := db.Get().Model(&entity.Playlist{})
	if style != nil {
		query = query.Where("style = ?", *style)
	}
	return query.Count(count).Error
}

func (p PlaylistRepo) GetPlaylistByTitle(playlist *entity.Playlist, title string) error {
	return db.Get().Model(&entity.Playlist{}).
		Where("title = ?", title).
		First(playlist).Error
}

func (p PlaylistRepo) CreatePlaylist(playlist *entity.Playlist) error {
	return db.Get().Create(playlist).Error
}

func (p PlaylistRepo) UpdatePlaylist(playlist *entity.Playlist, updateData *dto.PlaylistUpdateDTO) error {
	return db.Get().Model(&entity.Playlist{}).
		Where("id = ?", playlist.ID).
		Updates(updateData).Error
}

func (p PlaylistRepo) UpdatePlaylistCover(playlist *entity.Playlist, url string) error {
	return db.Get().Model(&entity.Playlist{}).
		Where("id = ?", playlist.ID).
		Update("cover_url", url).Error
}

func (p PlaylistRepo) GetPlaylistById(playlist *entity.Playlist, id uint64) error {
	return db.Get().Model(&entity.Playlist{}).
		Where("id = ?", id).
		First(playlist).Error
}

func (p PlaylistRepo) DeletePlaylist(playlist *entity.Playlist) error {
	return db.Get().Delete(playlist).Error
}

func (p PlaylistRepo) GetPlaylistCoverUrlsByIds(coverUrls *[]string, ids []uint64) error {
	return db.Get().Model(&entity.Playlist{}).
		Where("id IN ?", ids).
		Pluck("cover_url", coverUrls).Error
}

func (p PlaylistRepo) DeletePlaylists(ids []uint64) error {
	return db.Get().Delete(&entity.Playlist{}, ids).Error
}

func (p PlaylistRepo) GetAllPlaylistsByIds(data *result.PageResult[vo.PlaylistVO], userId uint64, ids []uint64, start int, size int, title *string, style *string) error {
	if len(ids) == 0 {
		data.Items = []vo.PlaylistVO{}
		data.Total = 0
		return nil
	}

	var vos []vo.PlaylistVO

	// 构建主查询
	query := db.Get().Table("tb_playlist p").
		Select("p.id AS playlist_id, p.title, p.cover_url").
		Joins("LEFT JOIN tb_user_favorite u ON p.id = u.playlist_id AND u.user_id = ?", userId)

	// WHERE p.id IN (...)
	query = query.Where("p.id IN ?", ids)

	// 模糊搜索 title
	if title != nil && *title != "" {
		query = query.Where("p.title LIKE ?", "%"+*title+"%")
	}

	// 模糊搜索 style（注意：原 SQL 用的是 LIKE，不是 =）
	if style != nil && *style != "" {
		query = query.Where("p.style LIKE ?", "%"+*style+"%")
	}

	// 排序：已收藏的按收藏时间倒序，未收藏的排后面
	// MySQL 中 NULL 会被排在最后（DESC 时），符合预期
	query = query.Order("u.create_time DESC")

	// 查询总数（用于分页）
	var total int64
	countDB := query.Session(&gorm.Session{NewDB: true}) // 避免污染原 query 的 SELECT
	if err := countDB.Select("COUNT(*)").Count(&total).Error; err != nil {
		return err
	}

	// 分页查询数据
	if err := query.Offset(start).Limit(size).Find(&vos).Error; err != nil {
		return err
	}

	data.Items = vos
	data.Total = total
	return nil
}
