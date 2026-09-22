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

// DonationService handles donations and usage records.
type DonationService struct {
	repo    *repository.DonationRepository
	usage   *repository.DonationUsageRepository
	orgRepo *repository.OrganizationRepository
	logger  *slog.Logger
}

// NewDonationService creates a DonationService.
func NewDonationService(repo *repository.DonationRepository, usage *repository.DonationUsageRepository, orgRepo *repository.OrganizationRepository, logger *slog.Logger) *DonationService {
	return &DonationService{repo: repo, usage: usage, orgRepo: orgRepo, logger: logger}
}

// Donate creates a donation (sandbox transaction id = mock).
func (s *DonationService) Donate(userID, orgID uint, amount float64) (*model.Donation, error) {
	if amount <= 0 {
		return nil, util.NewAppError(422, constants.CodeValidationError,
			fmt.Sprintf("Donation[org_id=%d] failed: amount must be positive", orgID))
	}
	org, err := s.orgRepo.FindByID(orgID)
	if err != nil {
		return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("Organization[id=%d] not found", orgID))
	}
	if org.Status != constants.OrgStatusApproved {
		return nil, util.NewAppError(403, constants.CodeForbidden, constants.MsgOrgNotApproved)
	}
	d := &model.Donation{
		UserID: userID, OrgID: orgID, Amount: amount,
		TransactionID: fmt.Sprintf("sandbox_%d", userID*100000+orgID), Status: "success",
	}
	if err := s.repo.Create(d); err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogDonationCreateFailed, orgID), "error", err)
		return nil, fmt.Errorf("donation create: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogDonationCreateSuccess, fmt.Sprintf("%.2f", amount), orgID), "id", d.ID)
	return d, nil
}

// ListByUser returns a user's donations.
func (s *DonationService) ListByUser(userID uint) ([]model.Donation, error) {
	items, err := s.repo.ListByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("donation list by user: %w", err)
	}
	return items, nil
}

// CreateUsage publishes a donation usage record (org side).
func (s *DonationService) CreateUsage(userID, donationID uint, amount float64, desc, proofURL string) (*model.DonationUsage, error) {
	d, err := s.repo.FindByID(donationID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("Donation[id=%d] not found", donationID))
		}
		return nil, fmt.Errorf("usage donation find: %w", err)
	}
	org, err := s.orgRepo.FindByUserID(userID)
	if err != nil || org.ID != d.OrgID {
		return nil, util.NewAppError(403, constants.CodeForbidden,
			fmt.Sprintf("DonationUsage[donation_id=%d] create failed: not org owner", donationID))
	}
	u := &model.DonationUsage{OrgID: d.OrgID, DonationID: donationID, Amount: amount, UsageDesc: desc, ProofURL: proofURL}
	if err := s.usage.Create(u); err != nil {
		return nil, fmt.Errorf("usage create: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogUsageCreateSuccess, d.OrgID), "id", u.ID)
	return u, nil
}

// ListUsageByOrg returns usage records for an org.
func (s *DonationService) ListUsageByOrg(orgID uint) ([]model.DonationUsage, error) {
	items, err := s.usage.ListByOrg(orgID)
	if err != nil {
		return nil, fmt.Errorf("usage list by org: %w", err)
	}
	return items, nil
}

// Stats returns donation totals for an org.
func (s *DonationService) Stats(orgID uint) (float64, error) {
	return s.repo.SumByOrg(orgID)
}
