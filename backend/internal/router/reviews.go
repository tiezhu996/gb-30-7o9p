package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gbadopt/gbadopt/internal/config"
	"github.com/gbadopt/gbadopt/internal/handler"
	"github.com/gbadopt/gbadopt/internal/middleware"
)

func registerReviewRoutes(v1 *gin.RouterGroup, cfg *config.Config, h *handler.ReviewHandler, limiter *middleware.RateLimiter) {
	reviews := v1.Group("/reviews", middleware.AuthRequired(cfg))
	reviews.GET("/me", h.ListMy)
	reviews.GET("/org", middleware.RequireRole("org"), h.ListOrg)
	reviews.POST("", middleware.RequireRole("org"), limiter.Limit(), h.Create)
	reviews.PUT("/:id/submit", h.Submit)
}
