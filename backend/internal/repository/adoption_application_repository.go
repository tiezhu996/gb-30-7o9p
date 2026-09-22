package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/gbadopt/gbadopt/internal/constants"
	"github.com/gbadopt/gbadopt/internal/model"
)

// AdoptionApplicationRepository handles application persistence.
type AdoptionApplicationRepository struct{ db *gorm.DB }

// NewAdoptionApplicationRepository creates the repository.
func NewAdoptionApplicationRepository(db *gorm.DB) *AdoptionApplicationRepository {
	return &AdoptionApplicationRepository{db: db}
}

// Create inserts an application.
func (r *AdoptionApplicationRepository) Create(a *model.AdoptionApplication) error {
	return translate(r.db.Create(a).Error)
}

// CreateTx inserts an application within an outer transaction.
func (r *AdoptionApplicationRepository) CreateTx(tx *gorm.DB, a *model.AdoptionApplication) error {
	return translate(tx.Create(a).Error)
}

// FindByID locates an application by id.
func (r *AdoptionApplicationRepository) FindByID(id uint) (*model.AdoptionApplication, error) {
	var a model.AdoptionApplication
	if err := translate(r.db.First(&a, id).Error); err != nil {
		return nil, err
	}
	return &a, nil
}

// FindByIDForUpdateTx locks an application row for the duration of tx.
func (r *AdoptionApplicationRepository) FindByIDForUpdateTx(tx *gorm.DB, id uint) (*model.AdoptionApplication, error) {
	var a model.AdoptionApplication
	if err := translate(tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&a, id).Error); err != nil {
		return nil, err
	}
	return &a, nil
}

// Update persists an application.
func (r *AdoptionApplicationRepository) Update(a *model.AdoptionApplication) error {
	return translate(r.db.Save(a).Error)
}

// UpdateTx persists an application within an outer transaction.
func (r *AdoptionApplicationRepository) UpdateTx(tx *gorm.DB, a *model.AdoptionApplication) error {
	return translate(tx.Save(a).Error)
}

// UpdateLifecycleTx writes lifecycle columns with an expected-status guard.
// It returns ErrConflict when another operation changed the row first, so a
// concurrent loser cannot overwrite the winner's record.
func (r *AdoptionApplicationRepository) UpdateLifecycleTx(tx *gorm.DB, a *model.AdoptionApplication, expectedStatus string) error {
	res := tx.Model(&model.AdoptionApplication{}).
		Where("id = ? AND status = ?", a.ID, expectedStatus).
		Updates(map[string]any{"status": a.Status, "end_reason": a.EndReason})
	if res.Error != nil {
		return translate(res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrConflict
	}
	return nil
}

// ListByUser returns applications of a user.
func (r *AdoptionApplicationRepository) ListByUser(userID uint) ([]model.AdoptionApplication, error) {
	var items []model.AdoptionApplication
	if err := r.db.Where("user_id = ?", userID).Order("id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// ListByOrg returns applications targeting an org.
func (r *AdoptionApplicationRepository) ListByOrg(orgID uint, status string) ([]model.AdoptionApplication, error) {
	var items []model.AdoptionApplication
	q := r.db.Where("org_id = ?", orgID)
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Order("pet_id ASC, created_at ASC, id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindByUserAndPet checks an existing application for the same pet.
func (r *AdoptionApplicationRepository) FindByUserAndPet(userID, petID uint) (*model.AdoptionApplication, error) {
	var a model.AdoptionApplication
	if err := translate(r.db.Where("user_id = ? AND pet_id = ?", userID, petID).First(&a).Error); err != nil {
		return nil, err
	}
	return &a, nil
}

// FindFirstWaitlistedForUpdateTx locks and returns the first active waitlist application.
func (r *AdoptionApplicationRepository) FindFirstWaitlistedForUpdateTx(tx *gorm.DB, petID uint) (*model.AdoptionApplication, error) {
	var a model.AdoptionApplication
	if err := translate(tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("pet_id = ? AND status = ?", petID, constants.AppStatusWaitlisted).
		Order("created_at ASC, id ASC").First(&a).Error); err != nil {
		return nil, err
	}
	return &a, nil
}

// ListActiveByPetIDs returns all applications that can still reserve or adopt the pets.
func (r *AdoptionApplicationRepository) ListActiveByPetIDs(petIDs []uint) ([]model.AdoptionApplication, error) {
	var items []model.AdoptionApplication
	if len(petIDs) == 0 {
		return items, nil
	}
	statuses := []string{
		constants.AppStatusSubmitted, constants.AppStatusOrgReview,
		constants.AppStatusCommunicating, constants.AppStatusConfirmed,
		constants.AppStatusOfflineInterview, constants.AppStatusReserved,
		constants.AppStatusWaitlisted,
	}
	if err := r.db.Where("pet_id IN ? AND status IN ?", petIDs, statuses).
		Order("pet_id ASC, created_at ASC, id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
