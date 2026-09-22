package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gbadopt/gbadopt/internal/config"
	"github.com/gbadopt/gbadopt/internal/handler"
	"github.com/gbadopt/gbadopt/internal/middleware"
)

func registerDonationRoutes(v1 *gin.RouterGroup, cfg *config.Config, h *handler.DonationHandler, limiter *middleware.RateLimiter) {
	donations := v1.Group("/donations", middleware.AuthRequired(cfg))
	donations.POST("", limiter.Limit(), h.Donate)
	donations.GET("/me", h.ListMy)
	donations.POST("/usage", middleware.RequireRole("org"), h.CreateUsage)
	v1.GET("/orgs/:id/usages", h.ListUsage)
	v1.GET("/orgs/:id/donation-stats", h.Stats)
}
