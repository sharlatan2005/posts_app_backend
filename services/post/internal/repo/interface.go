package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/sharlatan2005/posts_app_backend/services/post/internal/domain"
)

type PostRepo interface {
	Create(ctx context.Context, post *domain.Post) error
	Delete(ctx context.Context, postID uuid.UUID) error
	Update(ctx context.Context, postID uuid.UUID, newText string) error
	GetByID(ctx context.Context, postID uuid.UUID) (*domain.Post, error)

	GetAllUserPosts(ctx context.Context, userID uuid.UUID) ([]*domain.Post, error)
}
