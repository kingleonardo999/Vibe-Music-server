package repo

import (
	"vibe-music-server/internal/model/entity"
	"vibe-music-server/internal/model/vo"
	"vibe-music-server/internal/pkg/db"
	"vibe-music-server/internal/pkg/result"
)

type SongRepo struct{}

func NewSongRepo() *SongRepo {
	return &SongRepo{}
}

func (r SongRepo) GetAllSongs(data *result.PageResult[vo.SongVO], index, size int,
	songName, artistName, album *string) error {
	query := db.Get().Table("tb_song s").
		Select(`s.id            AS song_id,
		        s.name          AS song_name,
		        s.album,
		        s.duration,
		        s.cover_url     AS cover_url,
		        s.audio_url     AS audio_url,
		        s.release_time  AS release_time,
		        a.name          AS artist_name`).
		Joins("LEFT JOIN tb_artist a ON a.id = s.artist_id")

	// 动态条件
	if songName != nil {
		query = query.Where("s.name LIKE ?", "%"+*songName+"%")
	}
	if artistName != nil {
		query = query.Where("a.name LIKE ?", "%"+*artistName+"%")
	}
	if album != nil {
		query = query.Where("s.album LIKE ?", "%"+*album+"%")
	}

	// 总数
	if err := query.Count(&data.Total).Error; err != nil {
		return err
	}

	// 分页数据
	if err := query.
		Order("s.id DESC").
		Limit(size).
		Offset(index).
		Scan(&data.Items).Error; err != nil {
		return err
	}
	return nil
}

func (r SongRepo) GetAllSongsByIds(data *result.PageResult[vo.SongVO], ids []uint64, index, size int,
	songName, artistName, album *string) error {
	query := db.Get().Table("tb_song s").
		Select(`s.id            AS song_id,
		        s.name          AS song_name,
		        s.album,
		        s.duration,
		        s.cover_url     AS cover_url,
		        s.audio_url     AS audio_url,
		        s.release_time  AS release_time,
		        a.name          AS artist_name`).
		Joins("LEFT JOIN tb_artist a ON a.id = s.artist_id")

	// 动态条件
	if songName != nil {
		query = query.Where("s.name LIKE ?", "%"+*songName+"%")
	}
	if artistName != nil {
		query = query.Where("a.name LIKE ?", "%"+*artistName+"%")
	}
	if album != nil {
		query = query.Where("s.album LIKE ?", "%"+*album+"%")
	}

	query = query.Where("s.id IN ?", ids)

	// 总数
	if err := query.Count(&data.Total).Error; err != nil {
		return err
	}

	// 分页数据
	if err := query.
		Order("s.id DESC").
		Limit(size).
		Offset(index).
		Scan(&data.Items).Error; err != nil {
		return err
	}
	return nil
}

func (r SongRepo) GetAllSongsByArtist(data *result.PageResult[vo.SongAdminVO], name *string, album *string, id *uint64, index int, size int) error {
	query := db.Get().Table("tb_song s").
		Select(`s.id            AS song_id,
		        s.name          AS song_name,
		        s.album,
		        s.duration,
				s.style,
		        s.cover_url     AS cover_url,
		        s.audio_url     AS audio_url,
		        s.release_time  AS release_time,
		        a.name          AS artist_name`).
		Joins("LEFT JOIN tb_artist a ON a.id = s.artist_id")
	// 动态条件
	if name != nil {
		query = query.Where("a.name LIKE ?", "%"+*name+"%")
	}
	if album != nil {
		query = query.Where("s.album LIKE ?", "%"+*album+"%")
	}
	if id != nil {
		query = query.Where("a.id = ?", *id)
	}
	// 总数
	if err := query.Count(&data.Total).Error; err != nil {
		return err
	}
	// 分页数据
	if err := query.Order("release_time DESC").
		Limit(size).
		Offset(index).
		Scan(&data.Items).Error; err != nil {
		return err
	}
	return nil
}

func (r SongRepo) GetRandomSongs(data *[]vo.SongVO, limit int) error {
	query := db.Get().Table("tb_song s").
		Select(`s.id            AS song_id,
		        s.name          AS song_name,
		        s.album,
		        s.duration,
		        s.cover_url     AS cover_url,
		        s.audio_url     AS audio_url,
		        s.release_time  AS release_time,
		        a.name          AS artist_name`).
		Joins("LEFT JOIN tb_artist a ON a.id = s.artist_id").
		Order("RAND()").
		Limit(limit).
		Scan(data)
	return query.Error
}

func (r SongRepo) GetStylesByIds(data *[]string, ids []uint64) error {
	query := db.Get().Table("tb_song s").
		Select("s.style").
		Where("s.id IN ?", ids).
		Pluck("s.style", data)
	return query.Error
}

func (r SongRepo) GetRecommendedSongsByStyles(data *[]vo.SongVO, styles []string, ids []uint64, limit int) error {
	query := db.Get().Table("tb_song s").
		Select(`s.id            AS song_id,
		        s.name          AS song_name,
		        s.album,
		        s.duration,
		        s.cover_url     AS cover_url,
		        s.audio_url     AS audio_url,
		        s.release_time  AS release_time,
		        a.name          AS artist_name`).
		Joins("LEFT JOIN tb_artist a ON a.id = s.artist_id").
		Where("s.style IN ?", styles).
		Where("s.id NOT IN ?", ids).
		Order("RAND()").
		Limit(limit).
		Scan(data)
	return query.Error
}

func (r SongRepo) GetSongDetail(data *vo.SongDetailVO, id uint64) error {
	commentQuery := db.Get().Model(&entity.Comment{}).
		Where("type = ? AND playlist_id = ?", entity.CommentTypePlaylist, id).Scan(&data.Comments)
	query := db.Get().Table("tb_song s").
		Select(`s.id            AS song_id,
		        s.name          AS song_name,
		        s.album,
		        s.duration,
				s.style,
		        s.cover_url     AS cover_url,
		        s.audio_url     AS audio_url,
		        s.release_time  AS release_time,
		        a.name          AS artist_name`).
		Joins("LEFT JOIN tb_artist a ON a.id = s.artist_id").
		Where("s.id = ?", id).
		Scan(data)
	switch {
	case query.Error != nil:
		return query.Error
	case commentQuery.Error != nil:
		return commentQuery.Error
	}
	return nil
}

func (r SongRepo) GetAllSongsCount(count *int64, s *string) error {
	return db.Get().Model(&entity.Song{}).Where("name LIKE ?", "%"+*s+"%").Count(count).Error
}

func (r SongRepo) CreateSong(song *entity.Song) error {
	return db.Get().Create(song).Error
}

func (r SongRepo) GetSongById(song *entity.Song, id uint64) error {
	return db.Get().First(song, id).Error
}

func (r SongRepo) UpdateSong(song *entity.Song) error {
	return db.Get().Model(song).Updates(song).Error
}

func (r SongRepo) DeleteSongById(id uint64) error {
	return db.Get().Delete(&entity.Song{}, id).Error
}

func (r SongRepo) DeleteSongByIds(ids []uint64) error {
	return db.Get().Where("id IN ?", ids).Delete(&entity.Song{}).Error
}

func (r SongRepo) GetCoversByIds(covers *[]string, ids []uint64) error {
	return db.Get().Model(&entity.Song{}).
		Where("id IN ?", ids).
		Pluck("cover_url", covers).Error
}

func (r SongRepo) GetAudiosByIds(urls *[]string, ids []uint64) error {
	return db.Get().Model(&entity.Song{}).
		Where("id IN ?", ids).
		Pluck("audio_url", urls).Error
}
