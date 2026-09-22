package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gbadopt/gbadopt/internal/config"
	"github.com/gbadopt/gbadopt/internal/handler"
	"github.com/gbadopt/gbadopt/internal/middleware"
)

func registerUserRoutes(v1 *gin.RouterGroup, cfg *config.Config, h *handler.UserHandler, limiter *middleware.RateLimiter) {
	users := v1.Group("/users")
	users.POST("/register", limiter.Limit(), h.Register)
	users.POST("/login", limiter.Limit(), h.Login)
	me := users.Group("/me", middleware.AuthRequired(cfg))
	me.GET("", h.GetProfile)
	me.PUT("", h.UpdateProfile)
}
