package service

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/gbadopt/gbadopt/internal/constants"
	"github.com/gbadopt/gbadopt/internal/model"
	"github.com/gbadopt/gbadopt/internal/repository"
	"github.com/gbadopt/gbadopt/internal/util"
)

// PostService handles community posts and likes.
type PostService struct {
	repo   *repository.CommunityPostRepository
	orgRepo *repository.OrganizationRepository
	logger *slog.Logger
}

// NewPostService creates a PostService.
func NewPostService(repo *repository.CommunityPostRepository, orgRepo *repository.OrganizationRepository, logger *slog.Logger) *PostService {
	return &PostService{repo: repo, orgRepo: orgRepo, logger: logger}
}

// Create publishes a post (user or org).
func (s *PostService) Create(userID uint, p *model.CommunityPost) (*model.CommunityPost, error) {
	p.UserID = userID
	if p.Images == "" {
		p.Images = "[]"
	}
	if p.Status == "" {
		p.Status = "published"
	}
	if p.PostType == "" {
		p.PostType = "story"
	}
	if err := s.repo.Create(p); err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogPostCreateFailed, p.Title), "error", err)
		return nil, fmt.Errorf("post create: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogPostCreateSuccess, p.Title), "id", p.ID)
	return p, nil
}

// Get returns a post by id.
func (s *PostService) Get(id uint) (*model.CommunityPost, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("CommunityPost[id=%d] not found", id))
		}
		return nil, fmt.Errorf("post get: %w", err)
	}
	return p, nil
}

// Like increments a post's like count.
func (s *PostService) Like(id uint) (*model.CommunityPost, error) {
	if err := s.repo.IncrementLike(id); err != nil {
		return nil, fmt.Errorf("post like: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogPostLikeSuccess, id), "id", id)
	return s.repo.FindByID(id)
}

// List filters posts.
func (s *PostService) List(postType, keyword string, page, pageSize int) ([]model.CommunityPost, int64, error) {
	items, total, err := s.repo.List(postType, keyword, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("post list: %w", err)
	}
	return items, total, nil
}

// ListLatest returns recent posts.
func (s *PostService) ListLatest(limit int) ([]model.CommunityPost, error) {
	items, err := s.repo.ListLatest(limit)
	if err != nil {
		return nil, fmt.Errorf("post latest: %w", err)
	}
	return items, nil
}
