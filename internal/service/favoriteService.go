package service

import (
	"fmt"
	"time"
	"vibe-music-server/internal/model/dto"
	"vibe-music-server/internal/model/entity"
	"vibe-music-server/internal/model/vo"
	"vibe-music-server/internal/pkg/result"
	"vibe-music-server/internal/pkg/result/consts"
	"vibe-music-server/internal/pkg/util"
	"vibe-music-server/internal/repo"
)

type FavoriteService struct {
	favoriteRepo *repo.FavoriteRepo
	songRepo     *repo.SongRepo
	playlistRepo *repo.PlaylistRepo
}

func NewFavoriteService(favoriteRepo *repo.FavoriteRepo, songRepo *repo.SongRepo, playlistRepo *repo.PlaylistRepo) *FavoriteService {
	return &FavoriteService{
		favoriteRepo: favoriteRepo,
		songRepo:     songRepo,
		playlistRepo: playlistRepo,
	}
}

func (f FavoriteService) GetUserFavoriteSongs(songDTO *dto.SongDTO, claims *util.Claims) result.Result[result.PageResult[vo.SongVO]] {
	retErr := result.Error[result.PageResult[vo.SongVO]]
	retSuc := result.SuccessWithData[result.PageResult[vo.SongVO]]
	userId := claims.UserId
	pageNum := songDTO.PageNum
	pageSize := songDTO.PageSize
	start := (pageNum - 1) * pageSize
	var data result.PageResult[vo.SongVO]
	templateKey := fmt.Sprintf("favorite:getUserFavoriteSongs-%v-%v-%v-%v-%v", start, pageSize, songDTO.SongName, songDTO.ArtistName, songDTO.Album)
	if util.GetCache(templateKey, &data) {
		return retSuc(consts.Success, data)
	}
	var songIds []uint64
	if err := f.favoriteRepo.GetFavoriteSongIds(&songIds, userId); err != nil {
		return retErr(consts.InternalError)
	}
	if len(songIds) == 0 {
		return retSuc(consts.Success, result.PageResult[vo.SongVO]{
			Total: 0,
			Items: []vo.SongVO{},
		})
	}
	if err := f.songRepo.GetAllSongsByIds(&data, songIds, start, pageSize, songDTO.SongName, songDTO.ArtistName, songDTO.Album); err != nil {
		return retErr(consts.InternalError)
	}
	for _, song := range data.Items {
		song.LikeStatus = 1
	}
	util.SetCache(templateKey, data)
	return retSuc(consts.Success, data)
}

func (f FavoriteService) CollectSong(songId uint64, claims *util.Claims) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	userId := claims.UserId
	var isFavorite uint8
	if err := f.favoriteRepo.IsFavoriteSong(&isFavorite, userId, songId); err != nil {
		return retErr(consts.InternalError)
	}
	if isFavorite > 0 {
		return retErr(consts.Add + consts.Failed)
	}
	favorite := entity.Favorite{
		UserID:     userId,
		SongID:     &songId,
		Type:       entity.FavoriteTypeSong,
		CreateTime: time.Now(),
	}
	if err := f.favoriteRepo.AddFavorite(&favorite); err != nil {
		return retErr(consts.InternalError)
	}
	util.DeleteCacheByPattern("favorite:*")
	util.DeleteCacheByPattern("song:*")
	return retSuc(consts.Success)
}

func (f FavoriteService) CancelCollectSong(songId uint64, claims *util.Claims) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	userId := claims.UserId
	var isFavorite uint8
	if err := f.favoriteRepo.IsFavoriteSong(&isFavorite, userId, songId); err != nil {
		return retErr(consts.InternalError)
	}
	if isFavorite == 0 {
		return retErr(consts.Delete + consts.Failed)
	}
	if err := f.favoriteRepo.DeleteFavoriteSong(userId, songId); err != nil {
		return retErr(consts.InternalError)
	}
	util.DeleteCacheByPattern("favorite:*")
	util.DeleteCacheByPattern("song:*")
	return retSuc(consts.Success)
}

func (f FavoriteService) GetUserFavoritePlaylists(playlistDTO *dto.PlaylistDTO, claims *util.Claims) result.Result[result.PageResult[vo.PlaylistVO]] {
	retErr := result.Error[result.PageResult[vo.PlaylistVO]]
	retSuc := result.SuccessWithData[result.PageResult[vo.PlaylistVO]]
	userId := claims.UserId
	pageNum := playlistDTO.PageNum
	pageSize := playlistDTO.PageSize
	start := (pageNum - 1) * pageSize
	var data result.PageResult[vo.PlaylistVO]
	templateKey := fmt.Sprintf("favorite:getUserFavoritePlaylists-%v-%v-%v-%v", start, pageSize, playlistDTO.Title, playlistDTO.Style)
	if util.GetCache(templateKey, &data) {
		return retSuc(consts.Success, data)
	}
	var playlistIds []uint64
	if err := f.favoriteRepo.GetFavoritePlaylistIds(&playlistIds, userId); err != nil {
		return retErr(consts.InternalError)
	}
	if len(playlistIds) == 0 {
		return retSuc(consts.Success, result.PageResult[vo.PlaylistVO]{
			Total: 0,
			Items: []vo.PlaylistVO{},
		})
	}
	if err := f.playlistRepo.GetAllPlaylistsByIds(&data, userId, playlistIds, start, pageSize, playlistDTO.Title, playlistDTO.Style); err != nil {
		return retErr(consts.InternalError)
	}
	util.SetCache(templateKey, data)
	return retSuc(consts.Success, data)
}

func (f FavoriteService) CollectPlaylist(playlistId uint64, claims *util.Claims) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	userId := claims.UserId
	var isFavorite uint8
	if err := f.favoriteRepo.IsFavoritePlaylist(&isFavorite, userId, playlistId); err != nil {
		return retErr(consts.InternalError)
	}
	if isFavorite > 0 {
		return retErr(consts.Add + consts.Failed)
	}
	favorite := entity.Favorite{
		UserID:     userId,
		PlaylistID: &playlistId,
		Type:       entity.FavoriteTypePlaylist,
		CreateTime: time.Now(),
	}
	if err := f.favoriteRepo.AddFavorite(&favorite); err != nil {
		return retErr(consts.InternalError)
	}
	util.DeleteCacheByPattern("favorite:*")
	util.DeleteCacheByPattern("playlist:*")
	return retSuc(consts.Success)
}

func (f FavoriteService) CancelCollectPlaylist(playlistId uint64, claims *util.Claims) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	userId := claims.UserId
	var isFavorite uint8
	if err := f.favoriteRepo.IsFavoritePlaylist(&isFavorite, userId, playlistId); err != nil {
		return retErr(consts.InternalError)
	}
	if isFavorite == 0 {
		return retErr(consts.Delete + consts.Failed)
	}
	if err := f.favoriteRepo.DeleteFavoritePlaylist(userId, playlistId); err != nil {
		return retErr(consts.InternalError)
	}
	util.DeleteCacheByPattern("favorite:*")
	util.DeleteCacheByPattern("playlist:*")
	return retSuc(consts.Success)
}
