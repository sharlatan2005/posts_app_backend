package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/ctxutils"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/events"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/kafka/producer"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/repoerrors"
	"github.com/sharlatan2005/posts_app_backend/services/comment/internal/domain"
	"github.com/sharlatan2005/posts_app_backend/services/comment/internal/repo"
	"github.com/sharlatan2005/posts_app_backend/services/comment/internal/servErrors"
)

type CommentService struct {
	commentRepo   repo.CommentRepo
	kafkaProducer *producer.MyProducer
}

func NewCommentService(
	commentRepo repo.CommentRepo,
	kafkaProducer *producer.MyProducer) *CommentService {
	return &CommentService{
		commentRepo:   commentRepo,
		kafkaProducer: kafkaProducer,
	}
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

	go func(userID uuid.UUID) {

		activity := events.Activity{
			Type:         "comment",
			UserID:       userID,
			ActivityTime: time.Now(),
		}

		data, err := json.Marshal(activity)
		if err != nil {
			log.Printf("Ошибка сериализации: %v", err)
			return
		}

		err = s.kafkaProducer.SendMessage("activities", userID.String(), data)
		if err != nil {
			log.Printf("Ошибка отправки в Kafka: %v", err)
			return
		}
		log.Printf("Сообщение о выкладывании поста пользователя %s доставлено!", userID.String())

	}(comment.AuthorID)

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
	err := s.commentRepo.DeleteAllByPost(ctx, postID)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrNotFound):
			return servErrors.ErrNoCommentsOnPost
		default:
			return fmt.Errorf("Delete likes from post %s: %w", postID.String(), err)
		}
	}
	return nil
}
