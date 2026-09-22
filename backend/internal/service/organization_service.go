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

// OrganizationService handles org registration and certification.
type OrganizationService struct {
	repo   *repository.OrganizationRepository
	logger *slog.Logger
}

// NewOrganizationService creates an OrganizationService.
func NewOrganizationService(repo *repository.OrganizationRepository, logger *slog.Logger) *OrganizationService {
	return &OrganizationService{repo: repo, logger: logger}
}

// Register creates an org profile under a user account.
func (s *OrganizationService) Register(userID uint, o *model.Organization) (*model.Organization, error) {
	o.UserID = userID
	if o.Status == "" {
		o.Status = constants.OrgStatusPending
	}
	if err := s.repo.Create(o); err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogOrgRegisterFailed, o.Name), "error", err)
		return nil, fmt.Errorf("organization register: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogOrgRegisterSuccess, o.Name), "id", o.ID)
	return o, nil
}

// Get returns an org by id.
func (s *OrganizationService) Get(id uint) (*model.Organization, error) {
	o, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("Organization[id=%d] not found", id))
		}
		return nil, fmt.Errorf("organization get: %w", err)
	}
	return o, nil
}

// GetByUser returns the org owned by a user.
func (s *OrganizationService) GetByUser(userID uint) (*model.Organization, error) {
	o, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("organization get by user: %w", err)
	}
	return o, nil
}

// Review updates an org's certification status (admin only).
func (s *OrganizationService) Review(id uint, status string) (*model.Organization, error) {
	if !constants.IsValidOrganizationStatus(status) {
		return nil, util.NewAppError(422, constants.CodeValidationError,
			fmt.Sprintf("Organization[id=%d] review failed: invalid status %q", id, status))
	}
	o, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("organization review find: %w", err)
	}
	o.Status = status
	if err := s.repo.Update(o); err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogOrgReviewFailed, id), "error", err)
		return nil, fmt.Errorf("organization review update: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogOrgReviewSuccess, id, status), "id", id)
	return o, nil
}

// List filters orgs by status and keyword.
func (s *OrganizationService) List(status, keyword string, page, pageSize int) ([]model.Organization, int64, error) {
	items, total, err := s.repo.List(status, keyword, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("organization list: %w", err)
	}
	return items, total, nil
}

// ListApproved returns approved orgs.
func (s *OrganizationService) ListApproved(limit int) ([]model.Organization, error) {
	items, err := s.repo.ListApproved(limit)
	if err != nil {
		return nil, fmt.Errorf("organization approved list: %w", err)
	}
	return items, nil
}
