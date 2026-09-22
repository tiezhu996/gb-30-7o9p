package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/gbadopt/gbadopt/internal/model"
)

// VisitReviewRepository handles review persistence.
type VisitReviewRepository struct{ db *gorm.DB }

// NewVisitReviewRepository creates the repository.
func NewVisitReviewRepository(db *gorm.DB) *VisitReviewRepository { return &VisitReviewRepository{db: db} }

// Create inserts a review.
func (r *VisitReviewRepository) Create(v *model.VisitReview) error { return translate(r.db.Create(v).Error) }

// FindByID locates a review by id.
func (r *VisitReviewRepository) FindByID(id uint) (*model.VisitReview, error) {
	var v model.VisitReview
	if err := translate(r.db.First(&v, id).Error); err != nil {
		return nil, err
	}
	return &v, nil
}

// Update persists a review.
func (r *VisitReviewRepository) Update(v *model.VisitReview) error { return translate(r.db.Save(v).Error) }

// ListByUser returns reviews for a user.
func (r *VisitReviewRepository) ListByUser(userID uint) ([]model.VisitReview, error) {
	var items []model.VisitReview
	if err := r.db.Where("user_id = ?", userID).Order("due_date ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// ListByOrg returns reviews for an org.
func (r *VisitReviewRepository) ListByOrg(orgID uint) ([]model.VisitReview, error) {
	var items []model.VisitReview
	if err := r.db.Where("org_id = ?", orgID).Order("due_date ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// MarkOverdue flips pending reviews past due date to overdue.
func (r *VisitReviewRepository) MarkOverdue(userID uint) error {
	return r.db.Model(&model.VisitReview{}).
		Where("user_id = ? AND status = ? AND due_date < ?", userID, model.ReviewPending, time.Now()).
		Update("status", model.ReviewOverdue).Error
}
