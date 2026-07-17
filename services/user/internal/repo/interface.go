package repo

import (
	"context"

	"github.com/sharlatan2005/posts_app_backend/services/user/internal/domain"
)

type UserRepo interface {
	Create(ctx context.Context, user *domain.User) error
}
