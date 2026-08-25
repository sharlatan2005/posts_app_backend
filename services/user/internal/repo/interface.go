package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/sharlatan2005/posts_app_backend/services/user/internal/domain"
)

type UserRepo interface {
	Create(ctx context.Context, user *domain.User) error
	Exists(ctx context.Context, username string) (bool, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
}
