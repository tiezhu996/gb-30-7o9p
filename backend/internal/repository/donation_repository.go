package repository

import (
	"gorm.io/gorm"

	"github.com/gbadopt/gbadopt/internal/model"
)

// DonationRepository handles donation persistence.
type DonationRepository struct{ db *gorm.DB }

// NewDonationRepository creates the repository.
func NewDonationRepository(db *gorm.DB) *DonationRepository { return &DonationRepository{db: db} }

// Create inserts a donation.
func (r *DonationRepository) Create(d *model.Donation) error { return translate(r.db.Create(d).Error) }

// FindByID locates a donation by id.
func (r *DonationRepository) FindByID(id uint) (*model.Donation, error) {
	var d model.Donation
	if err := translate(r.db.First(&d, id).Error); err != nil {
		return nil, err
	}
	return &d, nil
}

// Update persists a donation.
func (r *DonationRepository) Update(d *model.Donation) error { return translate(r.db.Save(d).Error) }

// ListByUser returns donations made by a user.
func (r *DonationRepository) ListByUser(userID uint) ([]model.Donation, error) {
	var items []model.Donation
	if err := r.db.Where("user_id = ?", userID).Order("id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// ListByOrg returns donations received by an org.
func (r *DonationRepository) ListByOrg(orgID uint) ([]model.Donation, error) {
	var items []model.Donation
	if err := r.db.Where("org_id = ?", orgID).Order("id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// SumByOrg totals successful donations for an org.
func (r *DonationRepository) SumByOrg(orgID uint) (float64, error) {
	var sum float64
	if err := r.db.Model(&model.Donation{}).
		Where("org_id = ? AND status = ?", orgID, "success").
		Select("COALESCE(SUM(amount), 0)").Scan(&sum).Error; err != nil {
		return 0, err
	}
	return sum, nil
}
