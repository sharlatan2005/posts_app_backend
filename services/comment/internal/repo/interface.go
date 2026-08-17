package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/sharlatan2005/posts_app_backend/services/comment/internal/domain"
)

type CommentRepo interface {
	Create(ctx context.Context, comment *domain.Comment) error
	Delete(ctx context.Context, commentID uuid.UUID) error
	Update(ctx context.Context, commentID uuid.UUID, newText string) error
	GetByID(ctx context.Context, commentID uuid.UUID) (*domain.Comment, error)

	GetAllPostComments(ctx context.Context, postID uuid.UUID) ([]*domain.Comment, error)
	DeleteAllByPost(ctx context.Context, postID uuid.UUID) error
}
