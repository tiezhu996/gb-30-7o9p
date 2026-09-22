package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gbadopt/gbadopt/internal/config"
	"github.com/gbadopt/gbadopt/internal/handler"
	"github.com/gbadopt/gbadopt/internal/middleware"
)

func registerPostRoutes(v1 *gin.RouterGroup, cfg *config.Config, ph *handler.PostHandler, ch *handler.CommentHandler, limiter *middleware.RateLimiter) {
	posts := v1.Group("/posts")
	posts.GET("", ph.List)
	posts.GET("/:id", ph.Get)
	posts.GET("/:id/comments", ch.List)
	auth := posts.Group("", middleware.AuthRequired(cfg))
	auth.POST("", limiter.Limit(), ph.Create)
	auth.POST("/:id/comments", limiter.Limit(), ch.Create)
	auth.PUT("/:id/like", ph.Like)
	v1.DELETE("/comments/:id", middleware.AuthRequired(cfg), ch.Delete)
}
