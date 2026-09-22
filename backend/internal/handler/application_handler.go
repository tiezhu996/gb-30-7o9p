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

// ApplicationHandler exposes adoption application endpoints.
type ApplicationHandler struct {
	svc    *service.ApplicationService
	logger *slog.Logger
}

// NewApplicationHandler creates an ApplicationHandler.
func NewApplicationHandler(svc *service.ApplicationService, logger *slog.Logger) *ApplicationHandler {
	return &ApplicationHandler{svc: svc, logger: logger}
}

// Submit handles POST /applications.
func (h *ApplicationHandler) Submit(c *gin.Context) {
	var req dto.ApplicationSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam+": "+err.Error()))
		return
	}
	a, err := h.svc.Submit(middleware.GetUserID(c), req.PetID, req.Questionnaire)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(a))
}

// ListMy handles GET /applications/me.
func (h *ApplicationHandler) ListMy(c *gin.Context) {
	items, err := h.svc.ListByUser(middleware.GetUserID(c))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

// ListOrg handles GET /applications/org?status=.
func (h *ApplicationHandler) ListOrg(c *gin.Context) {
	status := c.Query("status")
	items, err := h.svc.ListByOrg(middleware.GetUserID(c), status)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

// UpdateStatus handles PUT /applications/:id/status.
func (h *ApplicationHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid application id"))
		return
	}
	var req dto.ApplicationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam))
		return
	}
	a, err := h.svc.UpdateStatus(middleware.GetUserID(c), uint(id), middleware.GetUserRole(c), req.Status)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(a))
}
