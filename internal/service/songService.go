package service

import (
	"fmt"
	"vibe-music-server/internal/model/dto"
	"vibe-music-server/internal/model/entity"
	"vibe-music-server/internal/model/vo"
	"vibe-music-server/internal/pkg/result"
	"vibe-music-server/internal/pkg/result/consts"
	"vibe-music-server/internal/pkg/util"
	"vibe-music-server/internal/repo"
)

type SongService struct {
	songRepo     *repo.SongRepo
	favoriteRepo *repo.FavoriteRepo
	styleRepo    *repo.StyleRepo
	genreRepo    *repo.GenreRepo
	minioService *MinioService
}

func NewSongService(songRepo *repo.SongRepo, favoriteRepo *repo.FavoriteRepo, styleRepo *repo.StyleRepo, genreRepo *repo.GenreRepo, minioService *MinioService) *SongService {
	return &SongService{
		songRepo:     songRepo,
		favoriteRepo: favoriteRepo,
		styleRepo:    styleRepo,
		genreRepo:    genreRepo,
		minioService: minioService,
	}
}

// GetAllSongs claims 可为nil
func (s SongService) GetAllSongs(songDTO *dto.SongDTO, claims *util.Claims) result.Result[result.PageResult[vo.SongVO]] {
	retErr := result.Error[result.PageResult[vo.SongVO]]
	retSuc := result.SuccessWithData[result.PageResult[vo.SongVO]]
	pageNum := songDTO.PageNum
	pageSize := songDTO.PageSize
	startIndex := (pageNum - 1) * pageSize
	var data result.PageResult[vo.SongVO]
	templateKey := fmt.Sprintf("song:getAllSongs:%v-%v-%v-%v-%v", startIndex, pageSize, *songDTO.SongName, *songDTO.ArtistName, *songDTO.Album)
	if !util.GetCache(templateKey, &data) {
		if err := s.songRepo.GetAllSongs(&data, startIndex, pageSize, songDTO.SongName, songDTO.ArtistName, songDTO.Album); err != nil {
			return retErr(consts.InternalError)
		}
		if data.Total == 0 {
			return retErr(consts.DataNotFound)
		}
		util.SetCache(templateKey, data)
	}
	if claims == nil {
		// 此时 LikeStatus 均为 0（默认）
		return retSuc(consts.Success, data)
	}
	if claims.Role == consts.User {
		userId := claims.UserId
		// 填充 LikeStatus
		var favoriteSongIds []uint64
		if err := s.favoriteRepo.GetFavoriteSongIds(&favoriteSongIds, userId); err != nil {
			return retErr(consts.InternalError)
		}
		for i := range data.Items {
			songId := data.Items[i].SongID
			if util.BinarySearch(favoriteSongIds, songId) != -1 {
				data.Items[i].LikeStatus = 1
			}
		}
	}
	return retSuc(consts.Success, data)
}

func (s SongService) GetAllSongsByArtist(songAndArtistDTO *dto.SongAndArtistDTO) result.Result[result.PageResult[vo.SongAdminVO]] {
	retErr := result.Error[result.PageResult[vo.SongAdminVO]]
	retSuc := result.SuccessWithData[result.PageResult[vo.SongAdminVO]]
	pageNum := songAndArtistDTO.PageNum
	pageSize := songAndArtistDTO.PageSize
	startIndex := (pageNum - 1) * pageSize
	var data result.PageResult[vo.SongAdminVO]
	templateKey := fmt.Sprintf("song:getAllSongsByArtist:%v-%v-%v-%v-%v", startIndex, pageSize, songAndArtistDTO.SongName, songAndArtistDTO.Album, songAndArtistDTO.ArtistID)
	if !util.GetCache(templateKey, &data) {
		if err := s.songRepo.GetAllSongsByArtist(&data, songAndArtistDTO.SongName, songAndArtistDTO.Album, songAndArtistDTO.ArtistID, startIndex, pageSize); err != nil {
			return retErr(consts.InternalError)
		}
		if data.Total == 0 {
			return retErr(consts.DataNotFound)
		}
		util.SetCache(templateKey, data)
	}
	return retSuc(consts.Success, data)
}

// GetRecommendedSongs claims 可为nil
func (s SongService) GetRecommendedSongs(claims *util.Claims) result.Result[[]vo.SongVO] {
	retErr := result.Error[[]vo.SongVO]
	retSuc := result.SuccessWithData[[]vo.SongVO]
	var data []vo.SongVO
	if claims == nil {
		// 未登录用户，随机推荐 10 首歌
		if err := s.songRepo.GetRandomSongs(&data, 10); err != nil {
			return retErr(consts.InternalError)
		}
		if len(data) == 0 {
			return retErr(consts.DataNotFound)
		}
		// 默认 LikeStatus 均为 0
		return retSuc(consts.Success, data)
	}
	if claims.Role == consts.User {
		userId := claims.UserId
		// 已登录用户，先获取用户喜欢的风格
		var favoriteSongIds []uint64
		if err := s.favoriteRepo.GetFavoriteSongIds(&favoriteSongIds, userId); err != nil {
			return retErr(consts.InternalError)
		}
		if len(favoriteSongIds) == 0 {
			// 用户没有喜欢的风格，随机推荐 10 首歌
			if err := s.songRepo.GetRandomSongs(&data, 10); err != nil {
				return retErr(consts.InternalError)
			}
			if len(data) == 0 {
				return retErr(consts.DataNotFound)
			}
			// 默认 LikeStatus 均为 0
			return retSuc(consts.Success, data)
		}
		// 根据用户喜欢的风格，推荐歌曲
		var favoriteStyles []string
		if err := s.songRepo.GetStylesByIds(&favoriteStyles, favoriteSongIds); err != nil {
			return retErr(consts.InternalError)
		}
		var favoriteStyleIds []uint64
		if err := s.styleRepo.GetStyleIdsByNames(&favoriteStyleIds, favoriteStyles); err != nil {
			return retErr(consts.InternalError)
		}
		var styleFrequency = make(map[uint64]int)
		for _, styleId := range favoriteStyleIds {
			styleFrequency[styleId]++
		}
		if err := s.songRepo.GetRecommendedSongsByStyles(&data, favoriteStyles, favoriteSongIds, 10); err != nil {
			return retErr(consts.InternalError)
		}
		for len(data) < 10 {
			var haveIds = make([]uint64, 0, len(data))
			for _, song := range data {
				haveIds = append(haveIds, song.SongID)
			}
			var supplement []vo.SongVO
			if err := s.songRepo.GetRandomSongs(&supplement, 10); err != nil {
				return retErr(consts.InternalError)
			}
			for _, song := range supplement {
				if util.BinarySearch(haveIds, song.SongID) == -1 {
					data = append(data, song)
					if len(data) >= 10 {
						break
					}
				}
			}
		}
	}
	return retSuc(consts.Success, data)
}

// GetSongDetail claims 可为nil
func (s SongService) GetSongDetail(songId uint64, claims *util.Claims) result.Result[vo.SongDetailVO] {
	retErr := result.Error[vo.SongDetailVO]
	retSuc := result.SuccessWithData[vo.SongDetailVO]
	var data vo.SongDetailVO
	templateKey := fmt.Sprintf("song:getSongDetail:%v", songId)
	if !util.GetCache(templateKey, &data) {
		if err := s.songRepo.GetSongDetail(&data, songId); err != nil {
			return retErr(consts.InternalError)
		}
		if data.SongID == 0 {
			return retErr(consts.DataNotFound)
		}
		util.SetCache(templateKey, data)
	}
	if claims == nil {
		return retSuc(consts.Success, data)
	}
	if claims.Role == consts.User {
		userId := claims.UserId
		var favoriteSongIds []uint64
		if err := s.favoriteRepo.GetFavoriteSongIds(&favoriteSongIds, userId); err != nil {
			return retErr(consts.InternalError)
		}
		if util.BinarySearch(favoriteSongIds, songId) != -1 {
			data.LikeStatus = 1
		}
	}
	return retSuc(consts.Success, data)
}

func (s SongService) GetAllSongsCount(style *string) result.Result[int64] {
	var count int64
	if err := s.songRepo.GetAllSongsCount(&count, style); err != nil {
		return result.Error[int64](consts.InternalError)
	}
	return result.SuccessWithData[int64](consts.Success, count)
}

func (s SongService) AddSong(songAddDTO *dto.SongAddDTO) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	song := entity.Song{
		Name:        songAddDTO.SongName,
		ArtistID:    uint(songAddDTO.ArtistID),
		Album:       songAddDTO.Album,
		Style:       songAddDTO.Style,
		ReleaseTime: songAddDTO.ReleaseTime,
	}
	if err := s.songRepo.CreateSong(&song); err != nil {
		return retErr(consts.Add + consts.Failed)
	}
	// 解析style
	styles := util.ParseStyle(songAddDTO.Style)
	var styleIds []uint64
	if err := s.styleRepo.GetStyleIdsByNames(&styleIds, styles); err != nil {
		return retErr(consts.InternalError)
	}
	for _, styleId := range styleIds {
		genre := entity.Genre{
			SongID:  uint64(song.ID),
			StyleID: styleId,
		}
		if err := s.genreRepo.CreateGenre(&genre); err != nil {
			return retErr(consts.Add + consts.Failed)
		}
	}
	util.DeleteCacheByPattern("song:*")
	return retSuc(consts.Add + consts.Success)
}

func (s SongService) UpdateSong(songUpdateDTO *dto.SongUpdateDTO) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	var song entity.Song
	if err := s.songRepo.GetSongById(&song, songUpdateDTO.SongID); err != nil {
		return retErr(consts.InternalError)
	}
	if song.ID == 0 {
		return retErr(consts.DataNotFound)
	}
	song = entity.Song{
		ID:          uint(songUpdateDTO.SongID),
		Name:        songUpdateDTO.SongName,
		ArtistID:    uint(songUpdateDTO.ArtistID),
		Album:       songUpdateDTO.Album,
		Style:       songUpdateDTO.Style,
		ReleaseTime: songUpdateDTO.ReleaseTime,
	}
	if err := s.songRepo.UpdateSong(&song); err != nil {
		return retErr(consts.Update + consts.Failed)
	}
	// 删除原有的风格关联
	if err := s.genreRepo.DeleteGenresBySongId(songUpdateDTO.SongID); err != nil {
		return retErr(consts.InternalError)
	}
	// 解析style
	styles := util.ParseStyle(songUpdateDTO.Style)
	var styleIds []uint64
	if err := s.styleRepo.GetStyleIdsByNames(&styleIds, styles); err != nil {
		return retErr(consts.InternalError)
	}
	for _, styleId := range styleIds {
		genre := entity.Genre{
			SongID:  songUpdateDTO.SongID,
			StyleID: styleId,
		}
		if err := s.genreRepo.CreateGenre(&genre); err != nil {
			return retErr(consts.Add + consts.Failed)
		}
	}
	util.DeleteCacheByPattern("song:*")
	return retSuc(consts.Update + consts.Success)
}

func (s SongService) UpdateSongCover(songId uint64, coverUrl string) result.Result[result.Nil] {
	var song entity.Song
	if err := s.songRepo.GetSongById(&song, songId); err != nil {
		return result.Error[result.Nil](consts.InternalError)
	}
	if song.ID == 0 {
		return result.Error[result.Nil](consts.DataNotFound)
	}
	song.CoverURL = coverUrl
	if err := s.songRepo.UpdateSong(&song); err != nil {
		return result.Error[result.Nil](consts.Update + consts.Failed)
	}
	util.DeleteCacheByPattern("song:*")
	return result.Success[result.Nil](consts.Update + consts.Success)
}

func (s SongService) UpdateSongAudio(songId uint64, audioUrl string) result.Result[result.Nil] {
	var song entity.Song
	if err := s.songRepo.GetSongById(&song, songId); err != nil {
		return result.Error[result.Nil](consts.InternalError)
	}
	if song.ID == 0 {
		return result.Error[result.Nil](consts.DataNotFound)
	}
	song.AudioURL = audioUrl
	if err := s.songRepo.UpdateSong(&song); err != nil {
		return result.Error[result.Nil](consts.Update + consts.Failed)
	}
	util.DeleteCacheByPattern("song:*")
	return result.Success[result.Nil](consts.Update + consts.Success)
}

func (s SongService) DeleteSong(songId uint64) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	var song entity.Song
	if err := s.songRepo.GetSongById(&song, songId); err != nil {
		return retErr(consts.InternalError)
	}
	if song.ID == 0 {
		return retErr(consts.DataNotFound)
	}
	// 删除歌曲文件和封面
	if err := s.minioService.DeleteFile(song.CoverURL); err != nil {
		return retErr(consts.Delete + consts.Failed)
	}
	if err := s.minioService.DeleteFile(song.AudioURL); err != nil {
		return retErr(consts.Delete + consts.Failed)
	}
	// 删除歌曲信息
	if err := s.songRepo.DeleteSongById(songId); err != nil {
		return retErr(consts.Delete + consts.Failed)
	}
	// 删除风格关联
	if err := s.genreRepo.DeleteGenresBySongId(songId); err != nil {
		return retErr(consts.InternalError)
	}
	util.DeleteCacheByPattern("song:*")
	return retSuc(consts.Delete + consts.Success)
}

func (s SongService) DeleteSongs(songIds []uint64) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	var covers []string
	var audios []string
	if err := s.songRepo.GetCoversByIds(&covers, songIds); err != nil {
		return retErr(consts.InternalError)
	}
	if err := s.songRepo.GetAudiosByIds(&audios, songIds); err != nil {
		return retErr(consts.InternalError)
	}
	for _, cover := range covers {
		if err := s.minioService.DeleteFile(cover); err != nil {
			return retErr(consts.Delete + consts.Failed)
		}
	}
	for _, audio := range audios {
		if err := s.minioService.DeleteFile(audio); err != nil {
			return retErr(consts.Delete + consts.Failed)
		}
	}
	if err := s.songRepo.DeleteSongByIds(songIds); err != nil {
		return retErr(consts.Delete + consts.Failed)
	}
	// 删除风格关联
	if err := s.genreRepo.DeleteGenresBySongIds(songIds); err != nil {
		return retErr(consts.InternalError)
	}
	util.DeleteCacheByPattern("song:*")
	return retSuc(consts.Delete + consts.Success)
}
