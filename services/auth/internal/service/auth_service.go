package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/auth"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/clients/errorsutils"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/clients/user"
	"github.com/sharlatan2005/posts_app_backend/services/auth/internal/serverrors"
)

type AuthService struct {
	jwtSecret  []byte
	userClient user.Client
}

func NewAuthService(secret string, userClient user.Client) *AuthService {
	return &AuthService{
		jwtSecret:  []byte(secret),
		userClient: userClient,
	}
}

type AuthResult struct {
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *AuthService) RegisterUser(ctx context.Context, username string, passwordHash string) (*AuthResult, error) {
	err := s.userClient.CreateUser(ctx, username, passwordHash)
	if err != nil {
		var extErr *errorsutils.ExternalServiceError
		switch {
		case errors.As(err, &extErr):
			return nil, extErr
		default:
			return nil, fmt.Errorf("user creation: %w", err)
		}
	}

	new_uuid := uuid.New()
	authStruct := auth.NewAuth(s.jwtSecret)
	token, err := authStruct.GenerateToken(new_uuid, username)

	if err != nil {
		return nil, fmt.Errorf("Failed to generate token: %w", err)
	}

	return &AuthResult{
		Token:     token,
		UserID:    new_uuid.String(),
		Username:  username,
		CreatedAt: time.Now(),
	}, nil
}

func (s *AuthService) LoginUser(ctx context.Context, username string, passwordHash string) (*AuthResult, error) {
	user, err := s.userClient.GetUserByUsername(ctx, username)
	if err != nil {
		var extErr *errorsutils.ExternalServiceError
		switch {
		case errors.As(err, &extErr):
			return nil, extErr
		default:
			return nil, fmt.Errorf("getting user by username: %w", err)
		}
	}

	if user.PasswordHash == passwordHash {
		authStruct := auth.NewAuth(s.jwtSecret)
		token, err := authStruct.GenerateToken(user.ID, username)
		if err != nil {
			return nil, fmt.Errorf("Failed to generate token: %w", err)
		}
		return &AuthResult{
			Token:     token,
			UserID:    user.ID.String(),
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		}, nil
	} else {
		return nil, serverrors.ErrWrongPassword
	}
}

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}
