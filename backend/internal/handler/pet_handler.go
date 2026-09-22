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

// PetHandler exposes pet endpoints.
type PetHandler struct {
	svc    *service.PetService
	logger *slog.Logger
}

// NewPetHandler creates a PetHandler.
func NewPetHandler(svc *service.PetService, logger *slog.Logger) *PetHandler {
	return &PetHandler{svc: svc, logger: logger}
}

// List handles GET /pets.
func (h *PetHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "12"))
	species := c.Query("species")
	status := c.Query("status")
	city := c.Query("city")
	keyword := c.Query("keyword")
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 12
	}
	items, total, err := h.svc.List(species, status, city, keyword, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.PageData{List: items, Total: total, Page: page, Size: pageSize}))
}

// Get handles GET /pets/:id.
func (h *PetHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid pet id"))
		return
	}
	p, err := h.svc.Get(uint(id))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(p))
}

// Publish handles POST /pets (org).
func (h *PetHandler) Publish(c *gin.Context) {
	var req dto.PetCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam+": "+err.Error()))
		return
	}
	p := &model.Pet{
		Name: req.Name, Species: req.Species, Breed: req.Breed, Age: req.Age,
		Gender: req.Gender, Size: req.Size, City: req.City, Description: req.Description,
		Personality: req.Personality, HealthStatus: req.HealthStatus,
		Neutered: req.Neutered, Vaccinated: req.Vaccinated, ImageURLs: req.ImageURLs,
	}
	created, err := h.svc.Publish(middleware.GetUserID(c), p)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(created))
}

// UpdateStatus handles PUT /pets/:id/status (org).
func (h *PetHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid pet id"))
		return
	}
	var req dto.PetStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam))
		return
	}
	p, err := h.svc.UpdateStatus(middleware.GetUserID(c), uint(id), req.Status)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(p))
}
