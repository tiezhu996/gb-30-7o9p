package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/gbadopt/gbadopt/internal/constants"
	"github.com/gbadopt/gbadopt/internal/dto"
	"github.com/gbadopt/gbadopt/internal/middleware"
	"github.com/gbadopt/gbadopt/internal/service"
	"github.com/gbadopt/gbadopt/internal/util"
)

// FavoriteHandler exposes favorite endpoints.
type FavoriteHandler struct {
	svc    *service.FavoriteService
	logger *slog.Logger
}

// NewFavoriteHandler creates a FavoriteHandler.
func NewFavoriteHandler(svc *service.FavoriteService, logger *slog.Logger) *FavoriteHandler {
	return &FavoriteHandler{svc: svc, logger: logger}
}

// List handles GET /favorites.
func (h *FavoriteHandler) List(c *gin.Context) {
	items, err := h.svc.ListByUser(middleware.GetUserID(c))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

// Add handles POST /favorites.
func (h *FavoriteHandler) Add(c *gin.Context) {
	var req dto.FavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam))
		return
	}
	f, err := h.svc.Add(middleware.GetUserID(c), req.TargetType, req.TargetID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(f))
}

// Remove handles DELETE /favorites/:targetType/:targetId.
func (h *FavoriteHandler) Remove(c *gin.Context) {
	targetType := c.Param("targetType")
	targetID, err := strconv.ParseUint(c.Param("targetId"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid target id"))
		return
	}
	if err := h.svc.Remove(middleware.GetUserID(c), targetType, uint(targetID)); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(gin.H{"removed": true}))
}
