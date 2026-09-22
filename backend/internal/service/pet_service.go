package service

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/gbadopt/gbadopt/internal/constants"
	"github.com/gbadopt/gbadopt/internal/model"
	"github.com/gbadopt/gbadopt/internal/repository"
	"github.com/gbadopt/gbadopt/internal/util"
)

// PetService handles pet publishing and management.
type PetService struct {
	repo         *repository.PetRepository
	orgRepo      *repository.OrganizationRepository
	redis        *util.RedisClient
	logger       *slog.Logger
}

// NewPetService creates a PetService.
func NewPetService(repo *repository.PetRepository, orgRepo *repository.OrganizationRepository, redis *util.RedisClient, logger *slog.Logger) *PetService {
	return &PetService{repo: repo, orgRepo: orgRepo, redis: redis, logger: logger}
}

// Publish creates a pet for an approved org.
func (s *PetService) Publish(userID uint, p *model.Pet) (*model.Pet, error) {
	org, err := s.orgRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(403, constants.CodeForbidden,
				fmt.Sprintf("Pet[publisher=%d] publish failed: org profile not found", userID))
		}
		return nil, fmt.Errorf("pet publish org find: %w", err)
	}
	if org.Status != constants.OrgStatusApproved {
		return nil, util.NewAppError(403, constants.CodeForbidden,
			fmt.Sprintf("Pet[org_id=%d] publish failed: %s", org.ID, constants.MsgOrgNotApproved))
	}
	if !constants.IsValidPetSpecies(p.Species) {
		return nil, util.NewAppError(422, constants.CodeValidationError,
			fmt.Sprintf("Pet[species=%s] publish failed: invalid species", p.Species))
	}
	if p.Status == "" {
		p.Status = constants.PetStatusAvailable
	}
	p.OrgID = org.ID
	if p.ImageURLs == "" {
		p.ImageURLs = "[]"
	}
	if err := s.repo.Create(p); err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogPetPublishFailed, p.Name), "error", err)
		return nil, fmt.Errorf("pet publish: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogPetPublishSuccess, p.Name, p.Species), "id", p.ID)
	return p, nil
}

// Get returns a pet by id.
func (s *PetService) Get(id uint) (*model.Pet, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("Pet[id=%d] not found", id))
		}
		return nil, fmt.Errorf("pet get: %w", err)
	}
	return p, nil
}

// UpdateStatus changes a pet's status (org owner or admin).
func (s *PetService) UpdateStatus(userID uint, petID uint, status string) (*model.Pet, error) {
	if !constants.IsValidPetStatus(status) {
		return nil, util.NewAppError(422, constants.CodeValidationError,
			fmt.Sprintf("Pet[id=%d] status=%s invalid", petID, status))
	}
	p, err := s.repo.FindByID(petID)
	if err != nil {
		return nil, fmt.Errorf("pet status find: %w", err)
	}
	if err := s.canManage(userID, p.OrgID); err != nil {
		return nil, err
	}
	p.Status = status
	if err := s.repo.Update(p); err != nil {
		return nil, fmt.Errorf("pet status update: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogPetStatusChanged, petID, status), "id", petID)
	return p, nil
}

// List filters pets.
func (s *PetService) List(species, status, city, keyword string, page, pageSize int) ([]model.Pet, int64, error) {
	items, total, err := s.repo.List(species, status, city, keyword, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("pet list: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogPetListSuccess, species, status, page), "total", total)
	return items, total, nil
}

// ListByOrg returns pets of an org.
func (s *PetService) ListByOrg(orgID uint) ([]model.Pet, error) {
	items, err := s.repo.ListByOrg(orgID, "")
	if err != nil {
		return nil, fmt.Errorf("pet list by org: %w", err)
	}
	return items, nil
}

// ListHot returns recent available pets.
func (s *PetService) ListHot(limit int) ([]model.Pet, error) {
	items, err := s.repo.ListHot(limit)
	if err != nil {
		return nil, fmt.Errorf("pet hot list: %w", err)
	}
	return items, nil
}

func (s *PetService) canManage(userID, orgID uint) error {
	org, err := s.orgRepo.FindByUserID(userID)
	if err != nil {
		return util.NewAppError(403, constants.CodeForbidden, "org profile not found")
	}
	if org.ID != orgID {
		return util.NewAppError(403, constants.CodeForbidden,
			fmt.Sprintf("Pet[org_id=%d] manage failed: user_id=%d is not org owner", orgID, userID))
	}
	return nil
}
