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

// FavoriteService handles pet/post favorites.
type FavoriteService struct {
	repo   *repository.FavoriteRepository
	logger *slog.Logger
}

// NewFavoriteService creates a FavoriteService.
func NewFavoriteService(repo *repository.FavoriteRepository, logger *slog.Logger) *FavoriteService {
	return &FavoriteService{repo: repo, logger: logger}
}

// Add favorites a target.
func (s *FavoriteService) Add(userID uint, targetType string, targetID uint) (*model.Favorite, error) {
	if targetType != "pet" && targetType != "post" {
		return nil, util.NewAppError(422, constants.CodeValidationError,
			fmt.Sprintf("Favorite[target_type=%s] add failed: invalid target type", targetType))
	}
	f := &model.Favorite{UserID: userID, TargetType: targetType, TargetID: targetID}
	if err := s.repo.Create(f); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, util.NewAppError(409, constants.CodeConflict,
				fmt.Sprintf("Favorite[user_id=%d target=%s:%d] add failed: already favorited", userID, targetType, targetID))
		}
		return nil, fmt.Errorf("favorite add: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogFavoriteAddSuccess, userID, targetType, targetID), "id", f.ID)
	return f, nil
}

// Remove deletes a favorite.
func (s *FavoriteService) Remove(userID uint, targetType string, targetID uint) error {
	f, err := s.repo.Find(userID, targetType, targetID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return util.NewAppError(404, constants.CodeNotFound,
				fmt.Sprintf("Favorite[user_id=%d target=%s:%d] not found", userID, targetType, targetID))
		}
		return fmt.Errorf("favorite remove find: %w", err)
	}
	if err := s.repo.Delete(f.ID); err != nil {
		return fmt.Errorf("favorite remove: %w", err)
	}
	return nil
}

// ListByUser returns a user's favorites.
func (s *FavoriteService) ListByUser(userID uint) ([]model.Favorite, error) {
	items, err := s.repo.ListByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("favorite list: %w", err)
	}
	return items, nil
}
