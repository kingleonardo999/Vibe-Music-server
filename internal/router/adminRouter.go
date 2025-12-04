package router

import (
	"github.com/gin-gonic/gin"
	"vibe-music-server/internal/controller"
	"vibe-music-server/internal/middleware"
)

func registerAdminRouter(r *gin.Engine, ctrl *controller.AdminCtrl) {
	g := r.Group("/admin")
	// 组级中间件
	// admin
	{
		g.POST("/register", ctrl.Register)
		g.POST("/login", ctrl.Login)
		g.POST("/logout", ctrl.Logout)
	}
	g.Use(middleware.AdminAuthMiddleware())
	// user management
	{
		g.GET("/getAllUsersCount", ctrl.GetAllUsersCount)
		g.POST("/getAllUsers", ctrl.GetAllUsers)
		g.POST("/addUser", ctrl.AddUser)
		g.PUT("/updateUser", ctrl.UpdateUser)
		g.PATCH("/updateUserStatus/:id/:status", ctrl.UpdateUserStatus)
		g.DELETE("/deleteUser/:id", ctrl.DeleteUser)
		g.DELETE("/deleteUsers", ctrl.DeleteUsers)
	}
	// artist management
	{
		g.GET("/getAllArtistsCount", ctrl.GetAllArtistsCount)
		g.GET("/getAllArtists", ctrl.GetAllArtists)
		g.POST("/addArtist", ctrl.AddArtist)
		g.PUT("/updateArtist", ctrl.UpdateArtist)
		g.PATCH("/updateArtistAvatar/:id", ctrl.UpdateArtistAvatar)
		g.DELETE("/deleteArtist/:id", ctrl.DeleteArtist)
		g.DELETE("/deleteArtists", ctrl.DeleteArtists)
	}
	// song management
	{
		g.GET("/getAllSongsCount", ctrl.GetAllSongsCount)
		g.GET("/getAllArtistNames", ctrl.GetAllArtistNames)
		g.POST("/getAllSongsByArtist", ctrl.GetAllSongsByArtist)
		g.POST("/addSong", ctrl.AddSong)
		g.PUT("/updateSong", ctrl.UpdateSong)
		g.PATCH("/updateSongCover/:id", ctrl.UpdateSongCover)
		g.PATCH("/updateSongAudio/:id", ctrl.UpdateSongAudio)
		g.DELETE("/deleteSong/:id", ctrl.DeleteSong)
		g.DELETE("/deleteSongs", ctrl.DeleteSongs)
	}
	// playlist management
	{
		g.GET("/getAllPlaylistsCount", ctrl.GetAllPlaylistsCount)
		g.POST("/getAllPlaylists", ctrl.GetAllPlaylists)
		g.POST("/addPlaylist", ctrl.AddPlaylist)
		g.PUT("/updatePlaylist", ctrl.UpdatePlaylist)
		g.PATCH("/updatePlaylistCover/:id", ctrl.UpdatePlaylistCover)
		g.DELETE("/deletePlaylist/:id", ctrl.DeletePlaylist)
		g.DELETE("/deletePlaylists", ctrl.DeletePlaylists)
	}
}
