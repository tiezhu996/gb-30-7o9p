package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gbadopt/gbadopt/internal/config"
	"github.com/gbadopt/gbadopt/internal/handler"
	"github.com/gbadopt/gbadopt/internal/middleware"
)

func registerApplicationRoutes(v1 *gin.RouterGroup, cfg *config.Config, h *handler.ApplicationHandler, limiter *middleware.RateLimiter) {
	apps := v1.Group("/applications", middleware.AuthRequired(cfg))
	apps.POST("", limiter.Limit(), h.Submit)
	apps.GET("/me", h.ListMy)
	apps.GET("/org", middleware.RequireRole("org"), h.ListOrg)
	apps.PUT("/:id/status", h.UpdateStatus)
}
