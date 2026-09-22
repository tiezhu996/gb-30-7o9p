package service

import (
	"errors"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"

	"github.com/gbadopt/gbadopt/internal/config"
	"github.com/gbadopt/gbadopt/internal/constants"
	"github.com/gbadopt/gbadopt/internal/model"
	"github.com/gbadopt/gbadopt/internal/repository"
	"github.com/gbadopt/gbadopt/internal/util"
)

// UserService handles registration, login and profile.
type UserService struct {
	repo   *repository.UserRepository
	logger *slog.Logger
	cfg    *config.Config
}

// NewUserService creates a UserService.
func NewUserService(repo *repository.UserRepository, logger *slog.Logger, cfg *config.Config) *UserService {
	return &UserService{repo: repo, logger: logger, cfg: cfg}
}

// Register creates a user and returns a JWT.
func (s *UserService) Register(username, email, password, nickname, phone string) (*model.User, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("user register: %w", err)
	}
	u := &model.User{Username: username, Email: email, PasswordHash: string(hash), Nickname: nickname, Phone: phone, Role: "user"}
	if u.Nickname == "" {
		u.Nickname = username
	}
	if err := s.repo.Create(u); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, "", util.NewAppError(409, constants.CodeConflict, constants.MsgUsernameTaken)
		}
		s.logger.Error(fmt.Sprintf(constants.LogUserRegisterFailed, username), "error", err)
		return nil, "", fmt.Errorf("user register: %w", err)
	}
	token, err := util.GenerateToken(u.ID, u.Username, u.Role, s.cfg.JWTSecret, s.cfg.JWTExpire)
	if err != nil {
		return nil, "", fmt.Errorf("user register token: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogUserRegisterSuccess, username), "user_id", u.ID)
	return u, token, nil
}

// Login verifies credentials and returns a JWT.
func (s *UserService) Login(identifier, password string) (*model.User, string, error) {
	u, err := s.repo.FindByUsername(identifier)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, "", util.NewAppError(401, constants.CodeUnauthorized, constants.MsgInvalidCredentials)
		}
		return nil, "", fmt.Errorf("user login find: %w", err)
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		s.logger.Warn(fmt.Sprintf(constants.LogUserLoginFailed, identifier))
		return nil, "", util.NewAppError(401, constants.CodeUnauthorized, constants.MsgInvalidCredentials)
	}
	token, err := util.GenerateToken(u.ID, u.Username, u.Role, s.cfg.JWTSecret, s.cfg.JWTExpire)
	if err != nil {
		return nil, "", fmt.Errorf("user login token: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogUserLoginSuccess, u.Username), "user_id", u.ID)
	return u, token, nil
}

// UpdateProfile updates a user's profile.
func (s *UserService) UpdateProfile(id uint, nickname, avatar, phone, bio string) (*model.User, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("user profile find: %w", err)
	}
	if nickname != "" {
		u.Nickname = nickname
	}
	if avatar != "" {
		u.Avatar = avatar
	}
	if phone != "" {
		u.Phone = phone
	}
	if err := s.repo.Update(u); err != nil {
		return nil, fmt.Errorf("user profile update: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogUserProfileUpdated, id), "user_id", id)
	return u, nil
}

// GetByID returns a user.
func (s *UserService) GetByID(id uint) (*model.User, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("user get: %w", err)
	}
	return u, nil
}
