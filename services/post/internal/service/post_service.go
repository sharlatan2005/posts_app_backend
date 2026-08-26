package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/clients/comment"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/clients/errorsutils"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/clients/like"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/clients/user"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/ctxutils"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/events"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/kafka/producer"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/repoerrors"
	"github.com/sharlatan2005/posts_app_backend/services/post/internal/domain"
	"github.com/sharlatan2005/posts_app_backend/services/post/internal/repo"
	"github.com/sharlatan2005/posts_app_backend/services/post/internal/servErrors"
)

type PostService struct {
	postRepo      repo.PostRepo
	userClient    user.Client
	commentClient comment.Client
	likeClient    like.Client
	kafkaProducer *producer.MyProducer
}

func NewPostService(postRepo repo.PostRepo,
	userClient user.Client,
	commentClient comment.Client,
	likeClient like.Client,
	kafkaProducer *producer.MyProducer) *PostService {
	return &PostService{
		postRepo:      postRepo,
		userClient:    userClient,
		commentClient: commentClient,
		likeClient:    likeClient,
		kafkaProducer: kafkaProducer,
	}
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

	go func(userID uuid.UUID) {

		activity := events.Activity{
			Type:         "post",
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

	}(post.AuthorID)

	return nil
}

func (s *PostService) Delete(ctx context.Context, postID uuid.UUID) error {

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrNotFound):
			return servErrors.ErrPostNotFound
		default:
			return err
		}
	}

	userID := ctxutils.GetUserID(ctx)
	if userID != post.AuthorID {
		return servErrors.ErrPostForbiddenAction
	}

	err = s.commentClient.DeleteCommentsByPost(ctx, postID)
	if err != nil {
		var extErr *errorsutils.ExternalServiceError
		switch {
		case errors.As(err, &extErr):
			if extErr.StatusCode != http.StatusNotFound {
				return extErr
			}
		default:
			return fmt.Errorf("Deleting comments on post: %w", err)
		}
	}

	err = s.likeClient.DeleteLikesByPost(ctx, postID)
	if err != nil {
		var extErr *errorsutils.ExternalServiceError
		switch {
		case errors.As(err, &extErr):
			if extErr.StatusCode != http.StatusNotFound {
				return extErr
			}
		default:
			return fmt.Errorf("Deleting likes on post: %w", err)
		}
	}

	err = s.postRepo.Delete(ctx, postID)
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

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrNotFound):
			return servErrors.ErrPostNotFound
		default:
			return err
		}
	}

	userID := ctxutils.GetUserID(ctx)
	if userID != post.AuthorID {
		return servErrors.ErrPostForbiddenAction
	}

	err = s.postRepo.Update(ctx, postID, newText)

	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrNotFound):
			return servErrors.ErrPostNotFound
		default:
			return fmt.Errorf("Update post %s: %w", postID.String(), err)
		}
	}
	return nil
}

func (s *PostService) GetAllUserPosts(ctx context.Context, username string) ([]*domain.Post, error) {
	u, err := s.userClient.GetUserByUsername(ctx, username)
	if err != nil {
		var extErr *errorsutils.ExternalServiceError
		switch {
		case errors.As(err, &extErr):
			return nil, extErr
		default:
			return nil, fmt.Errorf("getting user by username: %w", err)
		}
	}

	posts, err := s.postRepo.GetAllUserPosts(ctx, u.ID)
	if err != nil {
		switch {
		default:
			return nil, fmt.Errorf("getting all user posts: %w", err)
		}
	}

	return posts, nil
}
