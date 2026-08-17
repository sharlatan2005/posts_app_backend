package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/ctxutils"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/repoerrors"
	"github.com/sharlatan2005/posts_app_backend/services/like/internal/domain"
	"github.com/sharlatan2005/posts_app_backend/services/like/internal/repo"
	"github.com/sharlatan2005/posts_app_backend/services/like/internal/servErrors"
)

type LikeService struct {
	likeRepo repo.LikeRepo
}

func NewLikeService(likeRepo repo.LikeRepo) *LikeService {
	return &LikeService{likeRepo: likeRepo}
}

func (s *LikeService) AddLike(ctx context.Context, postID uuid.UUID) error {

	like := &domain.Like{
		ID:      uuid.New(),
		PostID:  postID,
		LikerID: ctxutils.GetUserID(ctx),
	}

	err := s.likeRepo.AddLike(ctx, like)
	if err != nil {
		switch {
		default:
			return fmt.Errorf("create like: %w", err)
		}
	}

	return nil
}

// Ошибки: ErrLikeNotFound
func (s *LikeService) RemoveLike(ctx context.Context, postID uuid.UUID) error {
	userID := ctxutils.GetUserID(ctx)

	err := s.likeRepo.RemoveLike(ctx, userID, postID)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrNotFound):
			return servErrors.ErrLikeNotFound
		default:
			return fmt.Errorf("Delete like from post %s: %w", postID.String(), err)
		}
	}
	return nil
}

func (s *LikeService) GetLikedUserIDs(ctx context.Context, postID uuid.UUID) ([]uuid.UUID, error) {
	likedUserIDs, err := s.likeRepo.GetLikedUserIDs(ctx, postID)
	if err != nil {
		switch {
		default:
			return nil, err
		}
	}

	return likedUserIDs, nil
}

func (s *LikeService) DeleteAllLikesByPost(ctx context.Context, postID uuid.UUID) error {
	err := s.likeRepo.DeleteAllByPost(ctx, postID)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrNotFound):
			return servErrors.ErrNoLikesOnPost
		default:
			return fmt.Errorf("Delete likes from post %s: %w", postID.String(), err)
		}
	}
	return nil
}
