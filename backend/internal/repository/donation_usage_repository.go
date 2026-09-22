package repository

import (
	"gorm.io/gorm"

	"github.com/gbadopt/gbadopt/internal/model"
)

// DonationUsageRepository handles usage-record persistence.
type DonationUsageRepository struct{ db *gorm.DB }

// NewDonationUsageRepository creates the repository.
func NewDonationUsageRepository(db *gorm.DB) *DonationUsageRepository {
	return &DonationUsageRepository{db: db}
}

// Create inserts a usage record.
func (r *DonationUsageRepository) Create(u *model.DonationUsage) error {
	return translate(r.db.Create(u).Error)
}

// ListByOrg returns usage records of an org.
func (r *DonationUsageRepository) ListByOrg(orgID uint) ([]model.DonationUsage, error) {
	var items []model.DonationUsage
	if err := r.db.Where("org_id = ?", orgID).Order("id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
