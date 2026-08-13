package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/sharlatan2005/posts_app_backend/services/post/internal/domain"
)

type PostRepo interface {
	Create(ctx context.Context, post *domain.Post) error
	Delete(ctx context.Context, post_id uuid.UUID) error
	Update(ctx context.Context, post_id uuid.UUID, newText string) error
}
