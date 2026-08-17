package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/ctxutils"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/repoerrors"
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

func (s *CommentService) Delete(ctx context.Context, commentID uuid.UUID) error {
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrNotFound):
			return servErrors.ErrCommentNotFound
		default:
			return err
		}
	}

	userID := ctxutils.GetUserID(ctx)
	if comment.AuthorID != userID {
		return servErrors.ErrCommentForbiddenAction
	}

	err = s.commentRepo.Delete(ctx, commentID)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrNotFound):
			return servErrors.ErrCommentNotFound
		default:
			return fmt.Errorf("Delete comment %s: %w", commentID.String(), err)
		}
	}
	return nil
}

func (s *CommentService) Update(ctx context.Context, commentID uuid.UUID, newText string) error {
	if newText == "" {
		return servErrors.ErrCommentTextEmpty
	}

	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrNotFound):
			return servErrors.ErrCommentNotFound
		default:
			return err
		}
	}

	userID := ctxutils.GetUserID(ctx)
	if comment.AuthorID != userID {
		return servErrors.ErrCommentForbiddenAction
	}

	err = s.commentRepo.Update(ctx, commentID, newText)

	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrNotFound):
			return servErrors.ErrCommentNotFound
		default:
			return fmt.Errorf("Update comment %s: %w", commentID.String(), err)
		}
	}
	return nil
}

func (s *CommentService) GetAllPostComments(ctx context.Context, postID uuid.UUID) ([]*domain.Comment, error) {
	posts, err := s.commentRepo.GetAllPostComments(ctx, postID)
	if err != nil {
		switch {
		default:
			return nil, fmt.Errorf("getting all user posts: %w", err)
		}
	}

	return posts, nil
}

func (s *CommentService) DeleteAllCommentsByPost(ctx context.Context, postID uuid.UUID) error {
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
