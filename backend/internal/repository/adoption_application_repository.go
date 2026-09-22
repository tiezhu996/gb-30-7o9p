package repository

import (
	"gorm.io/gorm"

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

// FindByIDForUpdateTx locates an application by id inside a transaction and
// takes a row lock so concurrent transitions cannot both succeed.
func (r *AdoptionApplicationRepository) FindByIDForUpdateTx(tx *gorm.DB, id uint) (*model.AdoptionApplication, error) {
	var a model.AdoptionApplication
	if err := translate(tx.Clauses(clauseLockingUpdate).First(&a, id).Error); err != nil {
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

// UpdateStatusIfTx conditionally changes status only while the current status
// is one of expectFrom. It is the idempotency/concurrency guard: a repeated or
// competing operation affects zero rows instead of overwriting a newer state.
// endReason is written only when non-empty; clearRank zeroes waitlist_rank.
func (r *AdoptionApplicationRepository) UpdateStatusIfTx(tx *gorm.DB, id uint, expectFrom []string, to, endReason string, clearRank bool) (int64, error) {
	updates := map[string]interface{}{
		"status":     to,
		"updated_at": gorm.Expr("NOW()"),
	}
	if endReason != "" {
		updates["end_reason"] = endReason
	}
	if clearRank {
		updates["waitlist_rank"] = 0
	}
	res := tx.Model(&model.AdoptionApplication{}).
		Where("id = ? AND status IN ?", id, expectFrom).
		Updates(updates)
	return res.RowsAffected, translate(res.Error)
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
	if err := q.Order("id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// ListByIDs returns applications matching the given ids.
func (r *AdoptionApplicationRepository) ListByIDs(ids []uint) ([]model.AdoptionApplication, error) {
	var items []model.AdoptionApplication
	if len(ids) == 0 {
		return items, nil
	}
	if err := r.db.Where("id IN ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindActiveByUserAndPet returns the user's in-progress application for a pet,
// or ErrNotFound when none exists (a terminated application may be re-filed).
func (r *AdoptionApplicationRepository) FindActiveByUserAndPet(userID, petID uint) (*model.AdoptionApplication, error) {
	var a model.AdoptionApplication
	err := translate(r.db.
		Where("user_id = ? AND pet_id = ? AND status IN ?", userID, petID, constants.ApplicationActiveStatuses()).
		First(&a).Error)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// ListActiveByPetForUpdateTx returns all in-progress applications of a pet
// ordered by submission time (id), locking the rows in the transaction.
func (r *AdoptionApplicationRepository) ListActiveByPetForUpdateTx(tx *gorm.DB, petID uint) ([]model.AdoptionApplication, error) {
	var items []model.AdoptionApplication
	err := tx.Clauses(clauseLockingUpdate).
		Where("pet_id = ? AND status IN ?", petID, constants.ApplicationActiveStatuses()).
		Order("created_at ASC, id ASC").
		Find(&items).Error
	return items, translate(err)
}

// ListWaitlistByPetForUpdateTx returns waitlisted applications ordered by rank,
// locking the rows in the transaction.
func (r *AdoptionApplicationRepository) ListWaitlistByPetForUpdateTx(tx *gorm.DB, petID uint) ([]model.AdoptionApplication, error) {
	var items []model.AdoptionApplication
	err := tx.Clauses(clauseLockingUpdate).
		Where("pet_id = ? AND status = ?", petID, constants.AppStatusWaitlisted).
		Order("waitlist_rank ASC, created_at ASC, id ASC").
		Find(&items).Error
	return items, translate(err)
}

// CloseOtherActiveTx ends every in-progress application of a pet except
// keepID with the given terminal status and end reason. Waitlist ranks are
// preserved so the final placement is visible. Returns rows affected.
func (r *AdoptionApplicationRepository) CloseOtherActiveTx(tx *gorm.DB, petID, keepID uint, terminal, endReason string) (int64, error) {
	res := tx.Model(&model.AdoptionApplication{}).
		Where("pet_id = ? AND id <> ? AND status IN ?", petID, keepID, constants.ApplicationActiveStatuses()).
		Updates(map[string]interface{}{
			"status":     terminal,
			"end_reason": endReason,
			"updated_at": gorm.Expr("NOW()"),
		})
	return res.RowsAffected, translate(res.Error)
}
