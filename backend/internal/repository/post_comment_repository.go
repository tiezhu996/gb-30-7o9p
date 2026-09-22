package repository

import (
	"gorm.io/gorm"

	"github.com/gbadopt/gbadopt/internal/model"
)

// PostCommentRepository handles comment persistence.
type PostCommentRepository struct{ db *gorm.DB }

// NewPostCommentRepository creates the repository.
func NewPostCommentRepository(db *gorm.DB) *PostCommentRepository {
	return &PostCommentRepository{db: db}
}

// Create inserts a comment.
func (r *PostCommentRepository) Create(c *model.PostComment) error {
	return translate(r.db.Create(c).Error)
}

// CreateTx inserts a comment within an outer transaction.
func (r *PostCommentRepository) CreateTx(tx *gorm.DB, c *model.PostComment) error {
	return translate(tx.Create(c).Error)
}

// FindByID locates a comment by id.
func (r *PostCommentRepository) FindByID(id uint) (*model.PostComment, error) {
	var c model.PostComment
	if err := translate(r.db.First(&c, id).Error); err != nil {
		return nil, err
	}
	return &c, nil
}

// Delete removes a comment by id.
func (r *PostCommentRepository) Delete(id uint) error {
	res := r.db.Delete(&model.PostComment{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// ListByPost returns comments for a post.
func (r *PostCommentRepository) ListByPost(postID uint) ([]model.PostComment, error) {
	var items []model.PostComment
	if err := r.db.Where("post_id = ?", postID).Order("id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
