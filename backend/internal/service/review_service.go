package service

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/gbadopt/gbadopt/internal/constants"
	"github.com/gbadopt/gbadopt/internal/model"
	"github.com/gbadopt/gbadopt/internal/repository"
	"github.com/gbadopt/gbadopt/internal/util"
)

// ReviewService handles adoption follow-up visits.
type ReviewService struct {
	repo    *repository.VisitReviewRepository
	appRepo *repository.AdoptionApplicationRepository
	orgRepo *repository.OrganizationRepository
	logger  *slog.Logger
}

// NewReviewService creates a ReviewService.
func NewReviewService(repo *repository.VisitReviewRepository, appRepo *repository.AdoptionApplicationRepository, orgRepo *repository.OrganizationRepository, logger *slog.Logger) *ReviewService {
	return &ReviewService{repo: repo, appRepo: appRepo, orgRepo: orgRepo, logger: logger}
}

// Create schedules a review for an approved application (org side).
func (s *ReviewService) Create(userID, applicationID uint, scheduledDays int) (*model.VisitReview, error) {
	app, err := s.appRepo.FindByID(applicationID)
	if err != nil {
		return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("AdoptionApplication[id=%d] not found", applicationID))
	}
	v := &model.VisitReview{
		ApplicationID: applicationID, UserID: app.UserID, OrgID: app.OrgID,
		ScheduledDays: scheduledDays, Status: model.ReviewPending, Photos: "[]",
	}
	if scheduledDays <= 0 {
		v.ScheduledDays = 30
	}
	v.DueDate = time.Now().AddDate(0, 0, v.ScheduledDays)
	if err := s.repo.Create(v); err != nil {
		return nil, fmt.Errorf("review create: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogReviewCreateSuccess, applicationID), "id", v.ID)
	return v, nil
}

// ListByUser lists reviews for a user, marking overdue first.
func (s *ReviewService) ListByUser(userID uint) ([]model.VisitReview, error) {
	if err := s.repo.MarkOverdue(userID); err != nil {
		s.logger.Warn("review overdue mark failed", "error", err)
	}
	items, err := s.repo.ListByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("review list by user: %w", err)
	}
	return items, nil
}

// ListByOrg lists reviews for an org owned by the user.
func (s *ReviewService) ListByOrg(userID uint) ([]model.VisitReview, error) {
	org, err := s.orgRepo.FindByUserID(userID)
	if err != nil {
		return nil, util.NewAppError(403, constants.CodeForbidden, "org profile not found")
	}
	items, err := s.repo.ListByOrg(org.ID)
	if err != nil {
		return nil, fmt.Errorf("review list by org: %w", err)
	}
	return items, nil
}

// Submit marks a review as submitted with photos and note.
func (s *ReviewService) Submit(userID, id uint, photos, note string) (*model.VisitReview, error) {
	v, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("VisitReview[id=%d] not found", id))
		}
		return nil, fmt.Errorf("review submit find: %w", err)
	}
	if v.UserID != userID {
		return nil, util.NewAppError(403, constants.CodeForbidden, fmt.Sprintf("VisitReview[id=%d] submit failed: not owner", id))
	}
	if photos == "" {
		photos = "[]"
	}
	v.Photos = photos
	v.Note = note
	v.Status = model.ReviewSubmitted
	if err := s.repo.Update(v); err != nil {
		return nil, fmt.Errorf("review submit update: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogReviewSubmitSuccess, id), "id", id)
	return v, nil
}
