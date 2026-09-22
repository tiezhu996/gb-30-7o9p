package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/gbadopt/gbadopt/internal/config"
	"github.com/gbadopt/gbadopt/internal/dto"
	"github.com/gbadopt/gbadopt/internal/handler"
	"github.com/gbadopt/gbadopt/internal/middleware"
	"github.com/gbadopt/gbadopt/internal/repository"
	"github.com/gbadopt/gbadopt/internal/service"
	"github.com/gbadopt/gbadopt/internal/util"
)

// Setup builds the gin engine.
func Setup(cfg *config.Config, db *gorm.DB, redis *util.RedisClient, minio *util.MinIOClient, logger *slog.Logger) *gin.Engine {
	userRepo := repository.NewUserRepository(db)
	orgRepo := repository.NewOrganizationRepository(db)
	petRepo := repository.NewPetRepository(db)
	appRepo := repository.NewAdoptionApplicationRepository(db)
	reviewRepo := repository.NewVisitReviewRepository(db)
	postRepo := repository.NewCommunityPostRepository(db)
	commentRepo := repository.NewPostCommentRepository(db)
	donationRepo := repository.NewDonationRepository(db)
	usageRepo := repository.NewDonationUsageRepository(db)
	favRepo := repository.NewFavoriteRepository(db)

	userService := service.NewUserService(userRepo, logger, cfg)
	orgService := service.NewOrganizationService(orgRepo, logger)
	petService := service.NewPetService(petRepo, orgRepo, redis, logger)
	appService := service.NewApplicationService(db, appRepo, petRepo, orgRepo, logger)
	reviewService := service.NewReviewService(reviewRepo, appRepo, orgRepo, logger)
	postService := service.NewPostService(postRepo, orgRepo, logger)
	commentService := service.NewCommentService(db, commentRepo, postRepo, logger)
	donationService := service.NewDonationService(donationRepo, usageRepo, orgRepo, logger)
	favService := service.NewFavoriteService(favRepo, logger)

	userHandler := handler.NewUserHandler(userService, logger)
	orgHandler := handler.NewOrganizationHandler(orgService, logger)
	petHandler := handler.NewPetHandler(petService, logger)
	appHandler := handler.NewApplicationHandler(appService, logger)
	reviewHandler := handler.NewReviewHandler(reviewService, logger)
	postHandler := handler.NewPostHandler(postService, logger)
	commentHandler := handler.NewCommentHandler(commentService, logger)
	donationHandler := handler.NewDonationHandler(donationService, logger)
	favHandler := handler.NewFavoriteHandler(favService, logger)
	uploadHandler := handler.NewUploadHandler(minio, logger)
	homeHandler := handler.NewHomeHandler(petService, postService, orgService, redis)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.CORS(cfg))
	r.Use(middleware.ErrorHandler(logger))

	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, dto.OK(gin.H{"status": "ok"})) })

	limiter := middleware.NewRateLimiter(cfg.RateLimitReq, cfg.RateLimitWin)
	v1 := r.Group("/api/v1")
	{
		v1.GET("/home/overview", homeHandler.Overview)
		v1.GET("/files/:key", uploadHandler.Get)
		registerUserRoutes(v1, cfg, userHandler, limiter)
		registerOrgRoutes(v1, cfg, orgHandler, limiter)
		registerPetRoutes(v1, cfg, petHandler, limiter)
		registerApplicationRoutes(v1, cfg, appHandler, limiter)
		registerReviewRoutes(v1, cfg, reviewHandler, limiter)
		registerPostRoutes(v1, cfg, postHandler, commentHandler, limiter)
		registerDonationRoutes(v1, cfg, donationHandler, limiter)
		registerFavoriteRoutes(v1, cfg, favHandler, limiter)
		v1.POST("/uploads", middleware.AuthRequired(cfg), limiter.Limit(), uploadHandler.Upload)
	}
	return r
}
