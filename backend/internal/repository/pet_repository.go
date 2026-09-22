package repository

import (
	"gorm.io/gorm"

	"github.com/gbadopt/gbadopt/internal/constants"
	"github.com/gbadopt/gbadopt/internal/model"
)

// PetRepository handles pet persistence.
type PetRepository struct{ db *gorm.DB }

// NewPetRepository creates a PetRepository.
func NewPetRepository(db *gorm.DB) *PetRepository { return &PetRepository{db: db} }

// Create inserts a pet.
func (r *PetRepository) Create(p *model.Pet) error { return translate(r.db.Create(p).Error) }

// FindByID locates a pet by id.
func (r *PetRepository) FindByID(id uint) (*model.Pet, error) {
	var p model.Pet
	if err := translate(r.db.First(&p, id).Error); err != nil {
		return nil, err
	}
	return &p, nil
}

// Update persists a pet.
func (r *PetRepository) Update(p *model.Pet) error { return translate(r.db.Save(p).Error) }

// UpdateTx persists a pet within an outer transaction.
func (r *PetRepository) UpdateTx(tx *gorm.DB, p *model.Pet) error { return translate(tx.Save(p).Error) }

// Delete removes a pet by id.
func (r *PetRepository) Delete(id uint) error {
	res := r.db.Delete(&model.Pet{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// List filters pets by species/status/city/breed/keyword with pagination.
func (r *PetRepository) List(species, status, city, keyword string, page, pageSize int) ([]model.Pet, int64, error) {
	var items []model.Pet
	var total int64
	q := r.db.Model(&model.Pet{})
	if species != "" {
		q = q.Where("species = ?", species)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if city != "" {
		q = q.Where("city = ?", city)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("name LIKE ? OR breed LIKE ? OR description LIKE ?", like, like, like)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// ListByOrg returns pets published by an org.
func (r *PetRepository) ListByOrg(orgID uint, status string) ([]model.Pet, error) {
	var items []model.Pet
	q := r.db.Where("org_id = ?", orgID)
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Order("id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// ListHot returns recent available pets for the home page.
func (r *PetRepository) ListHot(limit int) ([]model.Pet, error) {
	var items []model.Pet
	if err := r.db.Where("status = ?", constants.PetStatusAvailable).Order("id DESC").Limit(limit).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
