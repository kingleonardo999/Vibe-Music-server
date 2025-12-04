package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
	"vibe-music-server/internal/config"
	"vibe-music-server/internal/controller"
	"vibe-music-server/internal/middleware"
	"vibe-music-server/internal/repo"
	"vibe-music-server/internal/service"
)

var (
	adminRepo    *repo.AdminRepo
	artistRepo   *repo.ArtistRepo
	bannerRepo   *repo.BannerRepo
	commentRepo  *repo.CommentRepo
	favoriteRepo *repo.FavoriteRepo
	feedbackRepo *repo.FeedbackRepo
	genreRepo    *repo.GenreRepo
	playlistRepo *repo.PlaylistRepo
	songRepo     *repo.SongRepo
	styleRepo    *repo.StyleRepo
	userRepo     *repo.UserRepo
)

var (
	adminService    *service.AdminService
	artistService   *service.ArtistService
	bannerService   *service.BannerService
	commentService  *service.CommentService
	emailService    *service.EmailService
	favoriteService *service.FavoriteService
	feedbackService *service.FeedbackService
	minioService    *service.MinioService
	playlistService *service.PlaylistService
	songService     *service.SongService
	userService     *service.UserService
)

var (
	adminCtrl    *controller.AdminCtrl
	artistCtrl   *controller.ArtistCtrl
	bannerCtrl   *controller.BannerCtrl
	commentCtrl  *controller.CommentCtrl
	favoriteCtrl *controller.FavoriteCtrl
	feedbackCtrl *controller.FeedbackCtrl
	playlistCtrl *controller.PlaylistCtrl
	songCtrl     *controller.SongCtrl
	userCtrl     *controller.UserCtrl
)

func init() {
	adminRepo = repo.NewAdminRepo()
	artistRepo = repo.NewArtistRepo()
	bannerRepo = repo.NewBannerRepo()
	commentRepo = repo.NewCommentRepo()
	favoriteRepo = repo.NewFavoriteRepo()
	feedbackRepo = repo.NewFeedbackRepo()
	genreRepo = repo.NewGenreRepo()
	playlistRepo = repo.NewPlaylistRepo()
	songRepo = repo.NewSongRepo()
	styleRepo = repo.NewStyleRepo()
	userRepo = repo.NewUserRepo()
}

func init() {
	adminService = service.NewAdminService(adminRepo)
	artistService = service.NewArtistService(artistRepo, favoriteRepo, minioService)
	bannerService = service.NewBannerService(bannerRepo, minioService)
	commentService = service.NewCommentService(commentRepo)
	emailService = service.NewEmailService()
	favoriteService = service.NewFavoriteService(favoriteRepo, songRepo, playlistRepo)
	feedbackService = service.NewFeedbackService(feedbackRepo)
	minioService = service.NewMinioService()
	playlistService = service.NewPlaylistService(playlistRepo, favoriteRepo, styleRepo, minioService)
	songService = service.NewSongService(songRepo, favoriteRepo, styleRepo, genreRepo, minioService)
	userService = service.NewUserService(userRepo, emailService, minioService)
}

func init() {
	adminCtrl = controller.NewAdminCtrl(adminService, userService, artistService, songService, playlistService, minioService)
	artistCtrl = controller.NewArtistCtrl(artistService)
	bannerCtrl = controller.NewBannerCtrl(bannerService, minioService)
	commentCtrl = controller.NewCommentCtrl(commentService)
	favoriteCtrl = controller.NewFavoriteCtrl(favoriteService)
	feedbackCtrl = controller.NewFeedbackCtrl(feedbackService)
	playlistCtrl = controller.NewPlaylistCtrl(playlistService)
	songCtrl = controller.NewSongCtrl(songService)
	userCtrl = controller.NewUserCtrl(userService, minioService)
}

func setupCORS(corsCfg config.CORS) gin.HandlerFunc {
	//var compiledRegexes []*regexp.Regexp
	//for _, pattern := range corsCfg.AllowOrigins {
	//	// 检查是否包含正则表达式特征
	//	if (strings.HasPrefix(pattern, "^") && strings.HasSuffix(pattern, "$")) ||
	//		strings.Contains(pattern, "*") {
	//		if re, err := regexp.Compile(pattern); err == nil {
	//			compiledRegexes = append(compiledRegexes, re)
	//		}
	//	}
	//}
	//
	//exactOrigins := make(map[string]bool)
	//for _, origin := range corsCfg.AllowOrigins {
	//	if !strings.Contains(origin, "*") &&
	//		!(strings.HasPrefix(origin, "^") && strings.HasSuffix(origin, "$")) {
	//		exactOrigins[origin] = true
	//	}
	//}

	return cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
			//if origin == "" {
			//	return true
			//}
			//
			//// 精确匹配
			//if exactOrigins[origin] {
			//	return true
			//}
			//
			//// 正则匹配
			//for _, re := range compiledRegexes {
			//	if re.MatchString(origin) {
			//		return true
			//	}
			//}
			//
			//return false
		},
		AllowMethods:     corsCfg.AllowMethods,
		AllowHeaders:     corsCfg.AllowHeaders,
		ExposeHeaders:    corsCfg.ExposeHeaders,
		AllowCredentials: corsCfg.AllowCredentials,
		MaxAge:           time.Duration(corsCfg.MaxAge) * time.Second,
	})
}

func NewEngine() *gin.Engine {
	r := gin.Default()
	// 全局中间件
	r.Use(middleware.LoginMiddleware())
	r.Use(setupCORS(config.Get().App.CORS))
	// 业务分组
	registerAdminRouter(r, adminCtrl)
	registerArtistRouter(r, artistCtrl)
	registerBannerRouter(r, bannerCtrl)
	registerCommentRouter(r, commentCtrl)
	registerFavoriteRouter(r, favoriteCtrl)
	registerFeedbackRouter(r, feedbackCtrl)
	registerPlaylistRouter(r, playlistCtrl)
	registerSongRouter(r, songCtrl)
	registerUserRouter(r, userCtrl)
	return r
}
