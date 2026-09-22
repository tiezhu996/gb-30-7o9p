package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/gbadopt/gbadopt/internal/constants"
	"github.com/gbadopt/gbadopt/internal/dto"
	"github.com/gbadopt/gbadopt/internal/middleware"
	"github.com/gbadopt/gbadopt/internal/model"
	"github.com/gbadopt/gbadopt/internal/service"
	"github.com/gbadopt/gbadopt/internal/util"
)

// OrganizationHandler exposes organization endpoints.
type OrganizationHandler struct {
	svc    *service.OrganizationService
	logger *slog.Logger
}

// NewOrganizationHandler creates an OrganizationHandler.
func NewOrganizationHandler(svc *service.OrganizationService, logger *slog.Logger) *OrganizationHandler {
	return &OrganizationHandler{svc: svc, logger: logger}
}

// List handles GET /orgs.
func (h *OrganizationHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	keyword := c.Query("keyword")
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	items, total, err := h.svc.List(status, keyword, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.PageData{List: items, Total: total, Page: page, Size: pageSize}))
}

// Get handles GET /orgs/:id.
func (h *OrganizationHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid org id"))
		return
	}
	o, err := h.svc.Get(uint(id))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(o))
}

// Register handles POST /orgs (org role).
func (h *OrganizationHandler) Register(c *gin.Context) {
	var req dto.OrgRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam+": "+err.Error()))
		return
	}
	o := &model.Organization{
		Name: req.Name, LicenseURL: req.LicenseURL, CertType: req.CertType,
		Contact: req.Contact, City: req.City, Description: req.Description,
	}
	created, err := h.svc.Register(middleware.GetUserID(c), o)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(created))
}

// Review handles PUT /orgs/:id/review (admin).
func (h *OrganizationHandler) Review(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid org id"))
		return
	}
	var req dto.OrgReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam))
		return
	}
	o, err := h.svc.Review(uint(id), req.Status)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(o))
}
