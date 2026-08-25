package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/repoerrors"
	"github.com/sharlatan2005/posts_app_backend/services/user/internal/domain"
	"github.com/sharlatan2005/posts_app_backend/services/user/internal/repo"
	"github.com/sharlatan2005/posts_app_backend/services/user/internal/servErrors"
)

type UserService struct {
	userRepo repo.UserRepo
}

func NewUserService(userRepo repo.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

// Возвращает ошибки: ErrUsernameEmpty, ErrNameEmpty, ErrPasswordEmpty, ErrUserAlreadyExists
func (s *UserService) Create(ctx context.Context, username string, password_hash string) error {

	if username == "" {
		return servErrors.ErrUsernameEmpty
	}
	if password_hash == "" {
		return servErrors.ErrPasswordEmpty
	}

	user := &domain.User{
		ID:            uuid.New(),
		Username:      username,
		Password_hash: password_hash,
	}

	err := s.userRepo.Create(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrDuplicate):
			return servErrors.ErrUserAlreadyExists
		default:
			return fmt.Errorf("create user %s: %w", username, err)
		}
	}

	return nil
}

func (s *UserService) Exists(ctx context.Context, username string) (bool, error) {

	exists, err := s.userRepo.Exists(ctx, username)
	if err != nil {
		switch {
		default:
			return false, fmt.Errorf("Exists check %s: %w", username, err)
		}

	}
	return exists, nil
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrNotFound):
			return nil, servErrors.ErrUserNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrNotFound):
			return nil, servErrors.ErrUserNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}
