package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gbadopt/gbadopt/internal/config"
	"github.com/gbadopt/gbadopt/internal/handler"
	"github.com/gbadopt/gbadopt/internal/middleware"
)

func registerPetRoutes(v1 *gin.RouterGroup, cfg *config.Config, h *handler.PetHandler, limiter *middleware.RateLimiter) {
	pets := v1.Group("/pets")
	pets.GET("", h.List)
	pets.GET("/:id", h.Get)
	pets.POST("", middleware.AuthRequired(cfg), middleware.RequireRole("org"), limiter.Limit(), h.Publish)
	pets.PUT("/:id/status", middleware.AuthRequired(cfg), middleware.RequireRole("org"), h.UpdateStatus)
}
