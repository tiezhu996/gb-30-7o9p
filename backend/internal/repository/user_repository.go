package repository

import (
	"gorm.io/gorm"

	"github.com/gbadopt/gbadopt/internal/model"
)

// UserRepository handles user persistence.
type UserRepository struct{ db *gorm.DB }

// NewUserRepository creates a UserRepository.
func NewUserRepository(db *gorm.DB) *UserRepository { return &UserRepository{db: db} }

// Create inserts a user.
func (r *UserRepository) Create(u *model.User) error { return translate(r.db.Create(u).Error) }

// FindByUsername locates a user by username.
func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var u model.User
	if err := translate(r.db.Where("username = ?", username).First(&u).Error); err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByID locates a user by id.
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var u model.User
	if err := translate(r.db.First(&u, id).Error); err != nil {
		return nil, err
	}
	return &u, nil
}

// Update persists a user.
func (r *UserRepository) Update(u *model.User) error { return translate(r.db.Save(u).Error) }
