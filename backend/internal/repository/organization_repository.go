package repository

import (
	"gorm.io/gorm"

	"github.com/gbadopt/gbadopt/internal/model"
)

// OrganizationRepository handles org persistence.
type OrganizationRepository struct{ db *gorm.DB }

// NewOrganizationRepository creates an OrganizationRepository.
func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository { return &OrganizationRepository{db: db} }

// Create inserts an org.
func (r *OrganizationRepository) Create(o *model.Organization) error { return translate(r.db.Create(o).Error) }

// FindByID locates an org by id.
func (r *OrganizationRepository) FindByID(id uint) (*model.Organization, error) {
	var o model.Organization
	if err := translate(r.db.First(&o, id).Error); err != nil {
		return nil, err
	}
	return &o, nil
}

// FindByUserID locates an org owned by a user.
func (r *OrganizationRepository) FindByUserID(userID uint) (*model.Organization, error) {
	var o model.Organization
	if err := translate(r.db.Where("user_id = ?", userID).First(&o).Error); err != nil {
		return nil, err
	}
	return &o, nil
}

// Update persists an org.
func (r *OrganizationRepository) Update(o *model.Organization) error { return translate(r.db.Save(o).Error) }

// List filters orgs by status and keyword with pagination.
func (r *OrganizationRepository) List(status, keyword string, page, pageSize int) ([]model.Organization, int64, error) {
	var items []model.Organization
	var total int64
	q := r.db.Model(&model.Organization{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("name LIKE ? OR city LIKE ?", like, like)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// ListApproved returns approved orgs for public pages.
func (r *OrganizationRepository) ListApproved(limit int) ([]model.Organization, error) {
	var items []model.Organization
	if err := r.db.Where("status = ?", "approved").Order("id DESC").Limit(limit).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
