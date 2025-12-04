package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"vibe-music-server/internal/model/dto"
	"vibe-music-server/internal/model/entity"
	"vibe-music-server/internal/pkg/result"
	"vibe-music-server/internal/pkg/result/consts"
	"vibe-music-server/internal/service"
)

type AdminCtrl struct {
	adminService    *service.AdminService
	userService     *service.UserService
	artistService   *service.ArtistService
	songService     *service.SongService
	playlistService *service.PlaylistService
	minioService    *service.MinioService
}

func NewAdminCtrl(adminService *service.AdminService,
	userService *service.UserService, artistService *service.ArtistService,
	songService *service.SongService, playlistService *service.PlaylistService,
	minioService *service.MinioService) *AdminCtrl {
	return &AdminCtrl{
		adminService:    adminService,
		userService:     userService,
		artistService:   artistService,
		songService:     songService,
		playlistService: playlistService,
		minioService:    minioService,
	}
}

func (a *AdminCtrl) Register(c *gin.Context) {
	var adminDTO dto.AdminDTO
	if err := c.ShouldBindJSON(&adminDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.adminService.Register(&adminDTO))
}

func (a *AdminCtrl) Login(c *gin.Context) {
	var adminDTO dto.AdminDTO
	if err := c.ShouldBindJSON(&adminDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[string](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.adminService.Login(&adminDTO))
}

func (a *AdminCtrl) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.adminService.Logout(token))
}

func (a *AdminCtrl) GetAllUsersCount(c *gin.Context) {
	c.JSON(http.StatusOK, a.userService.GetAllUsersCount())
}

func (a *AdminCtrl) GetAllUsers(c *gin.Context) {
	var userSearchDTO dto.UserSearchDTO
	if err := c.ShouldBindJSON(&userSearchDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.userService.GetAllUsers(&userSearchDTO))
}

func (a *AdminCtrl) AddUser(c *gin.Context) {
	var userAddDTO dto.UserAddDTO
	if err := c.ShouldBindJSON(&userAddDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.userService.AddUser(&userAddDTO))
}

func (a *AdminCtrl) UpdateUser(c *gin.Context) {
	var userDTO dto.UserDTO
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.userService.UpdateUser(&userDTO))
}

func (a *AdminCtrl) UpdateUserStatus(c *gin.Context) {
	id := c.Param("id")
	status := c.Param("status")
	if id == "" || status == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	userStatus, err := strconv.ParseInt(status, 10, 32)
	if err != nil || (userStatus != 0 && userStatus != 1) {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.userService.UpdateUserStatus(userId, entity.UserStatus(userStatus)))
}

func (a *AdminCtrl) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.userService.DeleteUser(userId))
}

func (a *AdminCtrl) DeleteUsers(c *gin.Context) {
	var ids []uint64
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.userService.DeleteUsers(ids))
}

func (a *AdminCtrl) GetAllArtistsCount(c *gin.Context) {
	var gender *uint8
	if gStr := c.Query("gender"); gStr != "" {
		g, err := strconv.ParseUint(gStr, 10, 8)
		u8g := uint8(g)
		if err != nil || (g != 0 && g != 1) {
			c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
			return
		}
		gender = &u8g
	}
	var area *string
	if aStr := c.Query("area"); aStr != "" {
		area = &aStr
	}
	c.JSON(http.StatusOK, a.artistService.GetAllArtistsCount(gender, area))
}

func (a *AdminCtrl) GetAllArtists(c *gin.Context) {
	var artistDTO dto.ArtistDTO
	if err := c.ShouldBindJSON(&artistDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.artistService.GetAllArtists(&artistDTO))
}

func (a *AdminCtrl) AddArtist(c *gin.Context) {
	var artistAddDTO dto.ArtistAddDTO
	if err := c.ShouldBindJSON(&artistAddDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.artistService.AddArtist(&artistAddDTO))
}

func (a *AdminCtrl) UpdateArtist(c *gin.Context) {
	var artistUpdateDTO dto.ArtistUpdateDTO
	if err := c.ShouldBindJSON(&artistUpdateDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.artistService.UpdateArtist(&artistUpdateDTO))
}

func (a *AdminCtrl) UpdateArtistAvatar(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	artistId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	avatar, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	avatarUrl, err := a.minioService.UploadFile(avatar, "artists")
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.Error[result.Nil](consts.FileUpload+consts.Failed))
		return
	}
	c.JSON(http.StatusOK, a.artistService.UpdateArtistAvatar(artistId, avatarUrl))
}

func (a *AdminCtrl) DeleteArtist(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	artistId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.artistService.DeleteArtist(artistId))
}

func (a *AdminCtrl) DeleteArtists(c *gin.Context) {
	var ids []uint64
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.artistService.DeleteArtists(ids))
}

func (a *AdminCtrl) GetAllSongsCount(c *gin.Context) {
	var style *string
	if styleStr := c.Query("style"); styleStr != "" {
		style = &styleStr
	}
	c.JSON(http.StatusOK, a.songService.GetAllSongsCount(style))
}

func (a *AdminCtrl) GetAllArtistNames(c *gin.Context) {
	c.JSON(http.StatusOK, a.artistService.GetAllArtistNames())
}

func (a *AdminCtrl) GetAllSongsByArtist(c *gin.Context) {
	var songAndArtistDTO dto.SongAndArtistDTO
	if err := c.ShouldBindJSON(&songAndArtistDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.songService.GetAllSongsByArtist(&songAndArtistDTO))
}

func (a *AdminCtrl) AddSong(c *gin.Context) {
	var songAddDTO dto.SongAddDTO
	if err := c.ShouldBindJSON(&songAddDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.songService.AddSong(&songAddDTO))
}

func (a *AdminCtrl) UpdateSong(c *gin.Context) {
	var songUpdateDTO dto.SongUpdateDTO
	if err := c.ShouldBindJSON(&songUpdateDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.songService.UpdateSong(&songUpdateDTO))
}

func (a *AdminCtrl) UpdateSongCover(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	songId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	cover, err := c.FormFile("cover")
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	coverUrl, err := a.minioService.UploadFile(cover, "songCovers")
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.Error[result.Nil](consts.FileUpload+consts.Failed))
		return
	}
	c.JSON(http.StatusOK, a.songService.UpdateSongCover(songId, coverUrl))
}

func (a *AdminCtrl) UpdateSongAudio(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	songId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	audio, err := c.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	audioUrl, err := a.minioService.UploadFile(audio, "songs")
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.Error[result.Nil](consts.FileUpload+consts.Failed))
		return
	}
	c.JSON(http.StatusOK, a.songService.UpdateSongAudio(songId, audioUrl))
}

func (a *AdminCtrl) DeleteSong(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	songId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.songService.DeleteSong(songId))
}

func (a *AdminCtrl) DeleteSongs(c *gin.Context) {
	var ids []uint64
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.songService.DeleteSongs(ids))
}

func (a *AdminCtrl) GetAllPlaylistsCount(c *gin.Context) {
	var style *string
	if styleStr := c.Query("style"); styleStr != "" {
		style = &styleStr
	}
	c.JSON(http.StatusOK, a.playlistService.GetAllPlaylistsCount(style))
}

func (a *AdminCtrl) GetAllPlaylists(c *gin.Context) {
	var playlistDTO dto.PlaylistDTO
	if err := c.ShouldBindJSON(&playlistDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.playlistService.GetAllPlaylists(&playlistDTO))
}

func (a *AdminCtrl) AddPlaylist(c *gin.Context) {
	var playlistAddDTO dto.PlaylistAddDTO
	if err := c.ShouldBindJSON(&playlistAddDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.playlistService.AddPlaylist(&playlistAddDTO))
}

func (a *AdminCtrl) UpdatePlaylist(c *gin.Context) {
	var playlistUpdateDTO dto.PlaylistUpdateDTO
	if err := c.ShouldBindJSON(&playlistUpdateDTO); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.playlistService.UpdatePlaylist(&playlistUpdateDTO))
}

func (a *AdminCtrl) UpdatePlaylistCover(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	playlistId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	cover, err := c.FormFile("cover")
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	coverUrl, err := a.minioService.UploadFile(cover, "playlistCovers")
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.Error[result.Nil](consts.FileUpload+consts.Failed))
		return
	}
	c.JSON(http.StatusOK, a.playlistService.UpdatePlaylistCover(playlistId, coverUrl))
}

func (a *AdminCtrl) DeletePlaylist(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	playlistId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.playlistService.DeletePlaylist(playlistId))
}

func (a *AdminCtrl) DeletePlaylists(c *gin.Context) {
	var ids []uint64
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, result.Error[result.Nil](consts.InvalidParams))
		return
	}
	c.JSON(http.StatusOK, a.playlistService.DeletePlaylists(ids))
}
