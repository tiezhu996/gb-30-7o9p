package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gbadopt/gbadopt/internal/config"
	"github.com/gbadopt/gbadopt/internal/handler"
	"github.com/gbadopt/gbadopt/internal/middleware"
)

func registerOrgRoutes(v1 *gin.RouterGroup, cfg *config.Config, h *handler.OrganizationHandler, limiter *middleware.RateLimiter) {
	orgs := v1.Group("/orgs")
	orgs.GET("", h.List)
	orgs.GET("/:id", h.Get)
	orgs.POST("", middleware.AuthRequired(cfg), middleware.RequireRole("org"), limiter.Limit(), h.Register)
	orgs.PUT("/:id/review", middleware.AuthRequired(cfg), middleware.RequireRole("admin"), h.Review)
}
