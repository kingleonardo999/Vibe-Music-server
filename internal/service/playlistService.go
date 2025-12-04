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

type PlaylistService struct {
	playlistRepo *repo.PlaylistRepo
	favoriteRepo *repo.FavoriteRepo
	styleRepo    *repo.StyleRepo
	minioService *MinioService
}

func NewPlaylistService(playlistRepo *repo.PlaylistRepo, favoriteRepo *repo.FavoriteRepo, styleRepo *repo.StyleRepo, minioService *MinioService) *PlaylistService {
	return &PlaylistService{
		playlistRepo: playlistRepo,
		favoriteRepo: favoriteRepo,
		styleRepo:    styleRepo,
		minioService: minioService,
	}
}

func (p PlaylistService) GetAllPlaylists(playlistDTO *dto.PlaylistDTO) result.Result[result.PageResult[vo.PlaylistVO]] {
	retErr := result.Error[result.PageResult[vo.PlaylistVO]]
	retSuc := result.SuccessWithData[result.PageResult[vo.PlaylistVO]]
	pageNum := playlistDTO.PageNum
	pageSize := playlistDTO.PageSize
	startIndex := (pageNum - 1) * pageSize
	style := playlistDTO.Style
	title := playlistDTO.Title
	var data result.PageResult[vo.PlaylistVO]
	templateKey := fmt.Sprintf("playlist:getAllPlaylists:%v-%v-%v-%v", title, style, startIndex, pageSize)
	if util.GetCache(templateKey, &data) {
		return retSuc(consts.Success, data)
	}
	if err := p.playlistRepo.GetAllPlaylists(&data, title, style, startIndex, pageSize); err != nil {
		return retErr(consts.InternalError)
	}
	if data.Total == 0 {
		return retErr(consts.DataNotFound)
	}
	util.SetCache(templateKey, data)
	return retSuc(consts.Success, data)
}

func (p PlaylistService) GetAllPlaylistsInfo(playlistDTO *dto.PlaylistDTO) result.Result[result.PageResult[entity.Playlist]] {
	retErr := result.Error[result.PageResult[entity.Playlist]]
	retSuc := result.SuccessWithData[result.PageResult[entity.Playlist]]
	pageNum := playlistDTO.PageNum
	pageSize := playlistDTO.PageSize
	startIndex := (pageNum - 1) * pageSize
	style := playlistDTO.Style
	title := playlistDTO.Title
	var data result.PageResult[entity.Playlist]
	templateKey := fmt.Sprintf("playlist:getAllPlaylistsInfo:%v:%v:%v:%v", title, style, startIndex, pageSize)
	if util.GetCache(templateKey, &data) {
		return retSuc(consts.Success, data)
	}
	if err := p.playlistRepo.GetAllPlaylistsInfo(&data, title, style, startIndex, pageSize); err != nil {
		return retErr(consts.InternalError)
	}
	if data.Total == 0 {
		return retErr(consts.DataNotFound)
	}
	util.SetCache(templateKey, data)
	return retSuc(consts.Success, data)
}

// GetRecommendedPlaylists claims 可为nil
func (p PlaylistService) GetRecommendedPlaylists(claims *util.Claims) result.Result[[]vo.PlaylistVO] {
	retErr := result.Error[[]vo.PlaylistVO]
	retSuc := result.SuccessWithData[[]vo.PlaylistVO]
	var data []vo.PlaylistVO
	if claims == nil {
		// 用户未登录，返回默认推荐
		if err := p.playlistRepo.GetRandomPlaylists(&data, 10); err != nil {
			return retErr(consts.InternalError)
		}
		return retSuc(consts.Success, data)
	}
	userId := claims.UserId
	var favoritePlaylistIds []uint64
	if err := p.favoriteRepo.GetFavoritePlaylistIds(&favoritePlaylistIds, userId); err != nil {
		return retErr(consts.InternalError)
	}
	if len(favoritePlaylistIds) == 0 {
		// 用户未收藏任何歌单，返回默认推荐
		if err := p.playlistRepo.GetRandomPlaylists(&data, 10); err != nil {
			return retErr(consts.InternalError)
		}
		return retSuc(consts.Success, data)
	}
	var favoriteStyles []string
	if err := p.playlistRepo.GetStylesByIds(&favoriteStyles, favoritePlaylistIds); err != nil {
		return retErr(consts.InternalError)
	}
	var favoriteStyleIds []uint64
	if err := p.styleRepo.GetStyleIdsByNames(&favoriteStyleIds, favoriteStyles); err != nil {
		return retErr(consts.InternalError)
	}
	var styleFrequency = make(map[uint64]int)
	for _, styleId := range favoriteStyleIds {
		styleFrequency[styleId]++
	}
	if err := p.playlistRepo.GetRecommendedPlaylistsByStyles(&data, favoriteStyles, favoritePlaylistIds, 10); err != nil {
		return retErr(consts.InternalError)
	}
	for len(data) < 10 {
		var haveIds = make([]uint64, 0, len(data))
		for _, playlist := range data {
			haveIds = append(haveIds, playlist.PlaylistID)
		}
		var supplement []vo.PlaylistVO
		if err := p.playlistRepo.GetRandomPlaylists(&supplement, 10); err != nil {
			return retErr(consts.InternalError)
		}
		for _, playlist := range supplement {
			if util.BinarySearch(haveIds, playlist.PlaylistID) == -1 {
				data = append(data, playlist)
				if len(data) >= 10 {
					break
				}
			}
		}
	}
	return retSuc(consts.Success, data)
}

// GetPlaylistDetail claims 可为nil
func (p PlaylistService) GetPlaylistDetail(playlistId uint64, claims *util.Claims) result.Result[vo.PlaylistDetailVO] {
	retErr := result.Error[vo.PlaylistDetailVO]
	retSuc := result.SuccessWithData[vo.PlaylistDetailVO]
	var data vo.PlaylistDetailVO
	templateKey := fmt.Sprintf("playlist:getPlaylistDetail:%v", playlistId)
	if util.GetCache(templateKey, &data) {
		if claims != nil {
			userId := claims.UserId
			var isFavorite uint8
			if err := p.favoriteRepo.IsFavoritePlaylist(&isFavorite, userId, playlistId); err != nil {
				return retErr(consts.InternalError)
			}
			data.LikeStatus = isFavorite
		}
		return retSuc(consts.Success, data)
	}
	if err := p.playlistRepo.GetPlaylistDetail(&data, playlistId); err != nil {
		return retErr(consts.InternalError)
	}
	if data.PlaylistID == 0 {
		return retErr(consts.DataNotFound)
	}
	if claims != nil {
		userId := claims.UserId
		var isFavorite uint8
		if err := p.favoriteRepo.IsFavoritePlaylist(&isFavorite, userId, playlistId); err != nil {
			return retErr(consts.InternalError)
		}
		data.LikeStatus = isFavorite
	}
	util.SetCache(templateKey, data)
	return retSuc(consts.Success, data)
}

func (p PlaylistService) GetAllPlaylistsCount(style *string) result.Result[int64] {
	var count int64
	if err := p.playlistRepo.GetAllPlaylistsCount(&count, style); err != nil {
		return result.Error[int64](consts.InternalError)
	}
	return result.SuccessWithData[int64](consts.Success, count)
}

func (p PlaylistService) AddPlaylist(playlistAddDTO *dto.PlaylistAddDTO) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	var playlist entity.Playlist
	if err := p.playlistRepo.GetPlaylistByTitle(&playlist, playlistAddDTO.Title); err == nil {
		return retErr(consts.Playlist + consts.AlreadyExists)
	}
	playlist = entity.Playlist{
		Title:        playlistAddDTO.Title,
		Introduction: playlistAddDTO.Introduction,
		Style:        playlistAddDTO.Style,
	}
	if err := p.playlistRepo.CreatePlaylist(&playlist); err != nil {
		return retErr(consts.Add + consts.Failed)
	}
	util.DeleteCacheByPattern("playlist:*")
	return retSuc(consts.Add + consts.Success)
}

func (p PlaylistService) UpdatePlaylist(playlistUpdateDTO *dto.PlaylistUpdateDTO) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	var playlist entity.Playlist
	if err := p.playlistRepo.GetPlaylistByTitle(&playlist, playlistUpdateDTO.Title); err == nil && playlist.ID != uint(playlistUpdateDTO.PlaylistID) {
		return retErr(consts.Playlist + consts.AlreadyExists)
	}
	if err := p.playlistRepo.UpdatePlaylist(&playlist, playlistUpdateDTO); err != nil {
		return retErr(consts.Update + consts.Failed)
	}
	util.DeleteCacheByPattern("playlist:*")
	return retSuc(consts.Update + consts.Success)
}

func (p PlaylistService) UpdatePlaylistCover(playlistId uint64, coverUrl string) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	var playlist entity.Playlist
	if err := p.playlistRepo.GetPlaylistById(&playlist, playlistId); err != nil {
		return retErr(consts.Playlist + consts.NotFound)
	}
	// 删除旧封面
	if err := p.minioService.DeleteFile(playlist.CoverURL); err != nil {
		return retErr(consts.Update + consts.Failed)
	}
	if err := p.playlistRepo.UpdatePlaylistCover(&playlist, coverUrl); err != nil {
		return retErr(consts.Update + consts.Failed)
	}
	util.DeleteCacheByPattern("playlist:*")
	return retSuc(consts.Update + consts.Success)
}

func (p PlaylistService) DeletePlaylist(playlistId uint64) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	var playlist entity.Playlist
	if err := p.playlistRepo.GetPlaylistById(&playlist, playlistId); err != nil {
		return retErr(consts.Playlist + consts.NotFound)
	}
	if err := p.minioService.DeleteFile(playlist.CoverURL); err != nil {
		return retErr(consts.Delete + consts.Failed)
	}
	if err := p.playlistRepo.DeletePlaylist(&playlist); err != nil {
		return retErr(consts.Delete + consts.Failed)
	}
	util.DeleteCacheByPattern("playlist:*")
	return retSuc(consts.Delete + consts.Success)
}

func (p PlaylistService) DeletePlaylists(playlistIds []uint64) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	var coverUrls []string
	if err := p.playlistRepo.GetPlaylistCoverUrlsByIds(&coverUrls, playlistIds); err != nil {
		return retErr(consts.Playlist + consts.NotFound)
	}
	for _, coverUrl := range coverUrls {
		if err := p.minioService.DeleteFile(coverUrl); err != nil {
			return retErr(consts.Delete + consts.Failed)
		}
	}
	if err := p.playlistRepo.DeletePlaylists(playlistIds); err != nil {
		return retErr(consts.Delete + consts.Failed)
	}
	util.DeleteCacheByPattern("playlist:*")
	return retSuc(consts.Delete + consts.Success)
}
