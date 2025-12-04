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

type BannerService struct {
	bannerRepo   *repo.BannerRepo
	minioService *MinioService
}

func NewBannerService(bannerRepo *repo.BannerRepo, minioService *MinioService) *BannerService {
	return &BannerService{
		bannerRepo:   bannerRepo,
		minioService: minioService,
	}
}

func (b BannerService) GetAllBanners(bannerDTO *dto.BannerDTO) result.Result[result.PageResult[entity.Banner]] {
	retErr := result.Error[result.PageResult[entity.Banner]]
	retSuc := result.SuccessWithData[result.PageResult[entity.Banner]]
	pageNum := bannerDTO.PageNum
	pageSize := bannerDTO.PageSize
	startIndex := (pageNum - 1) * pageSize
	var pageRet result.PageResult[entity.Banner]
	templateKey := fmt.Sprintf("banner:getAllBanners:%v-%v-%v", *bannerDTO.Status, pageNum, pageSize)
	if util.GetCache(templateKey, &pageRet) {
		return retSuc(consts.Success, pageRet)
	}
	err := b.bannerRepo.GetPageBanners(&pageRet, bannerDTO.Status, startIndex, pageSize)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return retErr(consts.DataNotFound)
		}
		return retErr(consts.InternalError)
	}
	if pageRet.Total == 0 {
		return retErr(consts.DataNotFound)
	}
	util.SetCache(templateKey, pageRet)
	return retSuc(consts.Success, pageRet)
}

func (b BannerService) AddBanner(bannerUrl string) result.Result[result.Nil] {
	banner := entity.Banner{
		BannerURL: bannerUrl,
		Status:    entity.BannerStatusEnable,
	}
	if err := b.bannerRepo.CreateBanner(&banner); err != nil {
		return result.Error[result.Nil](consts.InternalError)
	}
	util.DeleteCacheByPattern("banner:*")
	return result.Success[result.Nil](consts.Success)
}

func (b BannerService) UpdateBanner(bannerId uint64, bannerUrl string) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	var banner entity.Banner
	if err := b.bannerRepo.GetBannerById(&banner, bannerId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return retErr(consts.DataNotFound)
		}
		return retErr(consts.InternalError)
	}
	// 删除旧图
	b.minioService.DeleteFile(banner.BannerURL)
	banner.BannerURL = bannerUrl
	if err := b.bannerRepo.UpdateBanner(&banner); err != nil {
		return retErr(consts.Update + consts.Failed)
	}
	util.DeleteCacheByPattern("banner:*")
	return retSuc(consts.Update + consts.Success)
}

func (b BannerService) UpdateBannerStatus(bannerId uint64, bannerStatus entity.BannerStatus) result.Result[result.Nil] {
	var retErr = result.Error[result.Nil]
	var retSuc = result.Success[result.Nil]
	var banner entity.Banner
	if err := b.bannerRepo.GetBannerById(&banner, bannerId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return retErr(consts.DataNotFound)
		}
		return retErr(consts.InternalError)
	}
	banner.Status = bannerStatus
	if err := b.bannerRepo.UpdateBanner(&banner); err != nil {
		return retErr(consts.Update + consts.Failed)
	}
	util.DeleteCacheByPattern("banner:*")
	return retSuc(consts.Update + consts.Success)
}

func (b BannerService) DeleteBanner(bannerId uint64) result.Result[result.Nil] {
	var retErr = result.Error[result.Nil]
	var retSuc = result.Success[result.Nil]
	var banner entity.Banner
	if err := b.bannerRepo.GetBannerById(&banner, bannerId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return retErr(consts.DataNotFound)
		}
		return retErr(consts.InternalError)
	}
	// 删除旧图
	b.minioService.DeleteFile(banner.BannerURL)
	if err := b.bannerRepo.DeleteBannerById(bannerId); err != nil {
		return retErr(consts.Delete + consts.Failed)
	}
	util.DeleteCacheByPattern("banner:*")
	return retSuc(consts.Delete + consts.Success)
}

func (b BannerService) DeleteBanners(bannerIds []uint64) result.Result[result.Nil] {
	var retErr = result.Error[result.Nil]
	var retSuc = result.Success[result.Nil]
	var banners []entity.Banner
	if err := b.bannerRepo.GetBannerByIds(&banners, bannerIds); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return retErr(consts.DataNotFound)
		}
		return retErr(consts.InternalError)
	}
	// 删除旧图
	for _, banner := range banners {
		b.minioService.DeleteFile(banner.BannerURL)
	}
	if err := b.bannerRepo.DeleteBannerByIds(bannerIds); err != nil {
		return retErr(consts.Delete + consts.Failed)
	}
	util.DeleteCacheByPattern("banner:*")
	return retSuc(consts.Delete + consts.Success)
}

func (b BannerService) GetBannerList() result.Result[[]vo.BannerVO] {
	var retErr = result.Error[[]vo.BannerVO]
	var retSuc = result.SuccessWithData[[]vo.BannerVO]
	var banners []vo.BannerVO
	templateKey := "banner:getBannerList"
	if util.GetCache(templateKey, &banners) {
		return retSuc(consts.Success, banners)
	}
	// 返回9个
	if err := b.bannerRepo.GetBannerList(&banners, 9); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return retErr(consts.DataNotFound)
		}
		return retErr(consts.InternalError)
	}
	return retSuc(consts.Success, banners)
}
