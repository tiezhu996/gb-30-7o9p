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

// ReviewHandler exposes visit review endpoints.
type ReviewHandler struct {
	svc    *service.ReviewService
	logger *slog.Logger
}

// NewReviewHandler creates a ReviewHandler.
func NewReviewHandler(svc *service.ReviewService, logger *slog.Logger) *ReviewHandler {
	return &ReviewHandler{svc: svc, logger: logger}
}

// ListMy handles GET /reviews/me.
func (h *ReviewHandler) ListMy(c *gin.Context) {
	items, err := h.svc.ListByUser(middleware.GetUserID(c))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

// ListOrg handles GET /reviews/org.
func (h *ReviewHandler) ListOrg(c *gin.Context) {
	items, err := h.svc.ListByOrg(middleware.GetUserID(c))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

// Create handles POST /reviews (org).
func (h *ReviewHandler) Create(c *gin.Context) {
	var req dto.ReviewCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam))
		return
	}
	v, err := h.svc.Create(middleware.GetUserID(c), req.ApplicationID, req.ScheduledDays)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(v))
}

// Submit handles PUT /reviews/:id/submit.
func (h *ReviewHandler) Submit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid review id"))
		return
	}
	var req dto.ReviewSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam))
		return
	}
	v, err := h.svc.Submit(middleware.GetUserID(c), uint(id), req.Photos, req.Note)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(v))
}
