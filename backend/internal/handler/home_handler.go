package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gbadopt/gbadopt/internal/dto"
	"github.com/gbadopt/gbadopt/internal/service"
	"github.com/gbadopt/gbadopt/internal/util"
)

// HomeHandler aggregates home page data, cached in Redis.
type HomeHandler struct {
	petService *service.PetService
	postService *service.PostService
	orgService *service.OrganizationService
	redis      *util.RedisClient
}

// NewHomeHandler creates a HomeHandler.
func NewHomeHandler(petService *service.PetService, postService *service.PostService, orgService *service.OrganizationService, redis *util.RedisClient) *HomeHandler {
	return &HomeHandler{petService: petService, postService: postService, orgService: orgService, redis: redis}
}

// Overview handles GET /home/overview.
func (h *HomeHandler) Overview(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	const cacheKey = "gbadopt:home:overview"
	if h.redis != nil {
		if cached, err := h.redis.GetString(ctx, cacheKey); err == nil && cached != "" {
			c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(cached))
			return
		}
	}
	hotPets, err := h.petService.ListHot(8)
	if err != nil {
		c.Error(err)
		return
	}
	latestPosts, err := h.postService.ListLatest(5)
	if err != nil {
		c.Error(err)
		return
	}
	orgs, err := h.orgService.ListApproved(4)
	if err != nil {
		c.Error(err)
		return
	}
	body := dto.OK(gin.H{"hot_pets": hotPets, "latest_posts": latestPosts, "orgs": orgs})
	if h.redis != nil {
		if b, err := json.Marshal(body); err == nil {
			_ = h.redis.SetString(ctx, cacheKey, string(b), 30*time.Second)
		}
	}
	c.JSON(http.StatusOK, body)
}
