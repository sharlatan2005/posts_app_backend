package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/ctxutils"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/repoerrors"
	"github.com/sharlatan2005/posts_app_backend/services/post/internal/domain"
	"github.com/sharlatan2005/posts_app_backend/services/post/internal/repo"
	"github.com/sharlatan2005/posts_app_backend/services/post/internal/servErrors"
)

type PostService struct {
	postRepo repo.PostRepo
}

func NewPostService(postRepo repo.PostRepo) *PostService {
	return &PostService{postRepo: postRepo}
}

func (s *PostService) Create(ctx context.Context, text string) error {

	if text == "" {
		return servErrors.ErrPostTextEmpty
	}

	post := &domain.Post{
		ID:       uuid.New(),
		AuthorID: ctxutils.GetUserID(ctx),
		Text:     text,
	}

	err := s.postRepo.Create(ctx, post)
	if err != nil {
		switch {
		default:
			return fmt.Errorf("create post: %w", err)
		}
	}

	return nil
}

func (s *PostService) Delete(ctx context.Context, postID uuid.UUID) error {
	err := s.postRepo.Delete(ctx, postID)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrNotFound):
			return servErrors.ErrPostNotFound
		default:
			return fmt.Errorf("Delete post %s: %w", postID.String(), err)
		}
	}
	return nil
}

func (s *PostService) Update(ctx context.Context, postID uuid.UUID, newText string) error {
	if newText == "" {
		return servErrors.ErrPostTextEmpty
	}

	err := s.postRepo.Update(ctx, postID, newText)

	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrNotFound):
			return servErrors.ErrPostNotFound
		default:
			return fmt.Errorf("Delete post %s: %w", postID.String(), err)
		}
	}
	return nil
}
