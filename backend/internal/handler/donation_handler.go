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

// DonationHandler exposes donation endpoints.
type DonationHandler struct {
	svc    *service.DonationService
	logger *slog.Logger
}

// NewDonationHandler creates a DonationHandler.
func NewDonationHandler(svc *service.DonationService, logger *slog.Logger) *DonationHandler {
	return &DonationHandler{svc: svc, logger: logger}
}

// Donate handles POST /donations.
func (h *DonationHandler) Donate(c *gin.Context) {
	var req dto.DonateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam+": "+err.Error()))
		return
	}
	d, err := h.svc.Donate(middleware.GetUserID(c), req.OrgID, req.Amount)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(d))
}

// ListMy handles GET /donations/me.
func (h *DonationHandler) ListMy(c *gin.Context) {
	items, err := h.svc.ListByUser(middleware.GetUserID(c))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

// CreateUsage handles POST /donations/usage (org).
func (h *DonationHandler) CreateUsage(c *gin.Context) {
	var req dto.UsageCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam+": "+err.Error()))
		return
	}
	u, err := h.svc.CreateUsage(middleware.GetUserID(c), req.DonationID, req.Amount, req.UsageDesc, req.ProofURL)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(u))
}

// ListUsage handles GET /orgs/:orgId/usages.
func (h *DonationHandler) ListUsage(c *gin.Context) {
	orgID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid org id"))
		return
	}
	items, err := h.svc.ListUsageByOrg(uint(orgID))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

// Stats handles GET /orgs/:orgId/donation-stats.
func (h *DonationHandler) Stats(c *gin.Context) {
	orgID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid org id"))
		return
	}
	sum, err := h.svc.Stats(uint(orgID))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(gin.H{"org_id": orgID, "total_donated": sum}))
}
