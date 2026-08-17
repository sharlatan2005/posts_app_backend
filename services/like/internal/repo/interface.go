package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/sharlatan2005/posts_app_backend/services/like/internal/domain"
)

type LikeRepo interface {
	AddLike(ctx context.Context, like *domain.Like) error
	RemoveLike(ctx context.Context, userID uuid.UUID, postID uuid.UUID) error
	GetLikedUserIDs(ctx context.Context, postID uuid.UUID) ([]uuid.UUID, error)

	DeleteAllByPost(ctx context.Context, postID uuid.UUID) error
}
