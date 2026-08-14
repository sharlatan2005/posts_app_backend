package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/ctxutils"
	"github.com/sharlatan2005/posts_app_backend/services/comment/internal/domain"
	"github.com/sharlatan2005/posts_app_backend/services/comment/internal/repo"
	"github.com/sharlatan2005/posts_app_backend/services/comment/internal/servErrors"
)

type CommentService struct {
	commentRepo repo.CommentRepo
}

func NewCommentService(commentRepo repo.CommentRepo) *CommentService {
	return &CommentService{commentRepo: commentRepo}
}

func (s *CommentService) Create(ctx context.Context, postID uuid.UUID, text string) error {

	if text == "" {
		return servErrors.ErrCommentTextEmpty
	}

	comment := &domain.Comment{
		ID:       uuid.New(),
		AuthorID: ctxutils.GetUserID(ctx),
		PostID:   postID,
		Text:     text,
	}

	err := s.commentRepo.Create(ctx, comment)
	if err != nil {
		switch {
		default:
			return fmt.Errorf("create comment: %w", err)
		}
	}

	return nil
}

// func (s *PostService) Delete(ctx context.Context, postID uuid.UUID) error {
// 	err := s.postRepo.Delete(ctx, postID)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, repoerrors.ErrNotFound):
// 			return servErrors.ErrPostNotFound
// 		default:
// 			return fmt.Errorf("Delete post %s: %w", postID.String(), err)
// 		}
// 	}
// 	return nil
// }

// func (s *PostService) Update(ctx context.Context, postID uuid.UUID, newText string) error {
// 	if newText == "" {
// 		return servErrors.ErrPostTextEmpty
// 	}

// 	err := s.postRepo.Update(ctx, postID, newText)

// 	if err != nil {
// 		switch {
// 		case errors.Is(err, repoerrors.ErrNotFound):
// 			return servErrors.ErrPostNotFound
// 		default:
// 			return fmt.Errorf("Delete post %s: %w", postID.String(), err)
// 		}
// 	}
// 	return nil
// }
