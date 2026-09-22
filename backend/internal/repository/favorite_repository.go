package repository

import (
	"gorm.io/gorm"

	"github.com/gbadopt/gbadopt/internal/model"
)

// FavoriteRepository handles favorite persistence.
type FavoriteRepository struct{ db *gorm.DB }

// NewFavoriteRepository creates the repository.
func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository { return &FavoriteRepository{db: db} }

// Create inserts a favorite.
func (r *FavoriteRepository) Create(f *model.Favorite) error { return translate(r.db.Create(f).Error) }

// Find locates a favorite by user and target.
func (r *FavoriteRepository) Find(userID uint, targetType string, targetID uint) (*model.Favorite, error) {
	var f model.Favorite
	if err := translate(r.db.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).First(&f).Error); err != nil {
		return nil, err
	}
	return &f, nil
}

// Delete removes a favorite by id.
func (r *FavoriteRepository) Delete(id uint) error { return r.db.Delete(&model.Favorite{}, id).Error }

// ListByUser returns favorites of a user.
func (r *FavoriteRepository) ListByUser(userID uint) ([]model.Favorite, error) {
	var items []model.Favorite
	if err := r.db.Where("user_id = ?", userID).Order("id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
