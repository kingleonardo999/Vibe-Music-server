package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"vibe-music-server/internal/model/dto"
	"vibe-music-server/internal/model/entity"
	"vibe-music-server/internal/model/vo"
	"vibe-music-server/internal/pkg/result"
	"vibe-music-server/internal/pkg/result/consts"
	"vibe-music-server/internal/pkg/util"
	"vibe-music-server/internal/repo"
)

type ArtistService struct {
	artistRepo   *repo.ArtistRepo
	favoriteRepo *repo.FavoriteRepo
	minioService *MinioService
}

func NewArtistService(artistRepo *repo.ArtistRepo, favoriteRepo *repo.FavoriteRepo, minioService *MinioService) *ArtistService {
	return &ArtistService{
		artistRepo:   artistRepo,
		favoriteRepo: favoriteRepo,
		minioService: minioService,
	}
}

func (a ArtistService) GetAllArtists(artistDTO *dto.ArtistDTO) result.Result[result.PageResult[vo.ArtistVO]] {
	retErr := result.Error[result.PageResult[vo.ArtistVO]]
	retSuc := result.SuccessWithData[result.PageResult[vo.ArtistVO]]
	pageNum := artistDTO.PageNum
	pageSize := artistDTO.PageSize
	startIndex := (pageNum - 1) * pageSize
	var pageRet result.PageResult[vo.ArtistVO]
	templateKey := util.GenKeyByPattern("artist:getAllArtists", artistDTO.ArtistName, artistDTO.Gender, artistDTO.Area, pageNum, pageSize)
	if util.GetCache(templateKey, &pageRet) {
		return retSuc(consts.Success, pageRet)
	}
	// 未命中缓存或反序列化失败
	err := a.artistRepo.GetPageArtistsVO(&pageRet, artistDTO.ArtistName, artistDTO.Gender, artistDTO.Area, startIndex, pageSize)
	if err != nil {
		return retErr(consts.InternalError)
	}
	if pageRet.Total == 0 {
		return retErr(consts.DataNotFound)
	}
	// 将结果存入缓存
	util.SetCache(templateKey, pageRet)
	return retSuc(consts.Success, pageRet)
}

func (a ArtistService) GetAllArtistsAndDetail(artistDTO *dto.ArtistDTO) result.Result[result.PageResult[entity.Artist]] {
	retErr := result.Error[result.PageResult[entity.Artist]]
	retSuc := result.SuccessWithData[result.PageResult[entity.Artist]]
	pageNum := artistDTO.PageNum
	pageSize := artistDTO.PageSize
	startIndex := (pageNum - 1) * pageSize
	var pageRet result.PageResult[entity.Artist]
	templateKey := fmt.Sprintf("artist:getAllArtistsAndDetail:%v-%v-%v-%v-%v", *artistDTO.ArtistName, *artistDTO.Gender, *artistDTO.Area, pageNum, pageSize)
	if util.GetCache(templateKey, &pageRet) {
		return retSuc(consts.Success, pageRet)
	}
	err := a.artistRepo.GetPageArtists(&pageRet, artistDTO.ArtistName, artistDTO.Gender, artistDTO.Area, startIndex, pageSize)
	if err != nil {
		return retErr(consts.InternalError)
	}
	if pageRet.Total == 0 {
		return retErr(consts.DataNotFound)
	}
	// 将结果存入缓存
	util.SetCache(templateKey, pageRet)
	return retSuc(consts.Success, pageRet)
}

func (a ArtistService) GetAllArtistNames() result.Result[[]vo.ArtistNameVO] {
	retErr := result.Error[[]vo.ArtistNameVO]
	retSuc := result.SuccessWithData[[]vo.ArtistNameVO]
	var data []vo.ArtistNameVO

	templateKey := "artist:getAllArtistNames"
	if util.GetCache(templateKey, &data) {
		return retSuc(consts.Success, data)
	}
	err := a.artistRepo.GetAllArtistsName(&data)
	if err != nil {
		return retErr(consts.InternalError)
	}
	if len(data) == 0 {
		return retErr(consts.DataNotFound)
	}
	util.SetCache(templateKey, data)
	return retSuc(consts.Success, data)
}

func (a ArtistService) GetRandomArtists() result.Result[[]vo.ArtistVO] {
	retErr := result.Error[[]vo.ArtistVO]
	retSuc := result.SuccessWithData[[]vo.ArtistVO]
	var data []vo.ArtistVO
	err := a.artistRepo.GetRandomArtists(&data, 10)
	if err != nil {
		return retErr(consts.InternalError)
	}
	if len(data) == 0 {
		return retErr(consts.DataNotFound)
	}
	return retSuc(consts.Success, data)
}

func (a ArtistService) GetArtistDetail(artistId uint64, claims *util.Claims) result.Result[vo.ArtistDetailVO] {
	retErr := result.Error[vo.ArtistDetailVO]
	retSuc := result.SuccessWithData[vo.ArtistDetailVO]
	var data vo.ArtistDetailVO
	templateKey := fmt.Sprintf("artist:getArtistDetail:%v", artistId)
	if !util.GetCache(templateKey, &data) {
		err := a.artistRepo.GetArtistDetail(&data, artistId)
		if err != nil {
			return retErr(consts.InternalError)
		}
	}
	if claims == nil {
		util.SetCache(templateKey, data)
		return retSuc(consts.Success, data)
	}
	// 根据 token 识别用户并设置 LikeStatus
	role := claims.Role
	if role != consts.UserRole {
		return retSuc(consts.Success, data)
	}
	userId := claims.UserId
	var favoriteRet []uint64
	err := a.favoriteRepo.GetFavoriteSongIds(&favoriteRet, userId)
	if err != nil {
		return retErr(consts.InternalError)
	}
	for i := range data.Songs {
		if util.BinarySearch[uint64](favoriteRet, data.Songs[i].SongID) != -1 {
			data.Songs[i].LikeStatus = 1
		}
	}
	util.SetCache(templateKey, data)
	return retSuc(consts.Success, data)
}

func (a ArtistService) GetAllArtistsCount(gender *uint8, area *string) result.Result[int64] {
	retErr := result.Error[int64]
	retSuc := result.SuccessWithData[int64]
	num, err := a.artistRepo.GetArtistsCount(gender, area)
	if err != nil {
		return retErr(consts.InternalError)
	}
	return retSuc(consts.Success, num)
}

func (a ArtistService) AddArtist(artistAddDTO *dto.ArtistAddDTO) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	if a.artistRepo.ExistArtistByName(artistAddDTO.ArtistName) {
		return retErr(consts.Artist + consts.AlreadyExists)
	}
	artist := entity.Artist{
		Name:         artistAddDTO.ArtistName,
		Gender:       artistAddDTO.Gender,
		Birth:        artistAddDTO.Birth,
		Area:         artistAddDTO.Area,
		Introduction: artistAddDTO.Introduction,
	}
	if err := a.artistRepo.CreateArtist(&artist); err != nil {
		return retErr(consts.Add + consts.Failed)
	}
	// 清除相关缓存
	util.DeleteCacheByPattern("artist:*")
	return retSuc(consts.Add + consts.Success)
}

func (a ArtistService) UpdateArtist(artistUpdateDTO *dto.ArtistUpdateDTO) result.Result[result.Nil] {
	var retErr = result.Error[result.Nil]
	var retSuc = result.Success[result.Nil]
	var artist entity.Artist
	if err := a.artistRepo.SelectByName(&artist, artistUpdateDTO.ArtistName); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return retErr(consts.Artist + consts.AlreadyExists)
	}
	if err := a.artistRepo.SelectById(&artist, artistUpdateDTO.ArtistID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return retErr(consts.DataNotFound)
		}
		return retErr(consts.InternalError)
	}
	if err := a.artistRepo.UpdateArtist(&artist, artistUpdateDTO); err != nil {
		return retErr(consts.Update + consts.Failed)
	}
	util.DeleteCacheByPattern("artist:*")
	return retSuc(consts.Update + consts.Success)
}

func (a ArtistService) UpdateArtistAvatar(artistId uint64, avatar string) result.Result[result.Nil] {
	var retErr = result.Error[result.Nil]
	var retSuc = result.Success[result.Nil]
	var artist entity.Artist
	if err := a.artistRepo.SelectById(&artist, artistId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return retErr(consts.DataNotFound)
		}
		return retErr(consts.InternalError)
	}
	artist.Avatar = avatar
	if err := a.artistRepo.UpdateArtist(&artist, map[string]any{"avatar": avatar}); err != nil {
		return retErr(consts.Update + consts.Failed)
	}
	util.DeleteCacheByPattern("artist:*")
	return retSuc(consts.Update + consts.Success)
}

func (a ArtistService) DeleteArtist(artistId uint64) result.Result[result.Nil] {
	// 1. 查询歌手信息，获取头像 url
	var retErr = result.Error[result.Nil]
	var retSuc = result.Success[result.Nil]
	var artist entity.Artist
	if err := a.artistRepo.SelectById(&artist, artistId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return retErr(consts.DataNotFound)
		}
		return retErr(consts.InternalError)
	}
	avatarURL := artist.Avatar
	// 2. 先删除 Minio 的文件
	if avatarURL != "" {
		a.minioService.DeleteFile(avatarURL)
	}
	// 3. 删除数据库记录
	if err := a.artistRepo.DeleteArtistById(artistId); err != nil {
		return retErr(consts.Delete + consts.Failed)
	}
	// 4. 清除相关缓存
	util.DeleteCacheByPattern("artist:*")
	return retSuc(consts.Delete + consts.Success)
}

func (a ArtistService) DeleteArtists(artistIds []uint64) result.Result[result.Nil] {
	// 1. 查询歌手信息，获取头像 url
	var retErr = result.Error[result.Nil]
	var retSuc = result.Success[result.Nil]
	var avatars []string
	if err := a.artistRepo.GetAvatarsByIds(&avatars, artistIds); err != nil {
		return retErr(consts.InternalError)
	}
	// 2. 先删除 Minio 的文件
	if len(avatars) > 0 {
		for _, avatarURL := range avatars {
			if avatarURL != "" {
				a.minioService.DeleteFile(avatarURL)
			}
		}
	}
	// 3. 删除数据库记录
	if err := a.artistRepo.DeleteArtistsByIds(artistIds); err != nil {
		return retErr(consts.Delete + consts.Failed)
	}
	// 4. 清除相关缓存
	util.DeleteCacheByPattern("artist:*")
	return retSuc(consts.Delete + consts.Success)
}
