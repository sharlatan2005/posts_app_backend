package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
	token, err := s.generateToken(new_uuid, username)

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
		token, err := s.generateToken(user.ID, username)
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

func (s *AuthService) generateToken(userID uuid.UUID, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) validateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	}) // Вот эта хуерга под капотом проверяет подпись с помощью jwt секрета + проверяет,
	// не истек ли срок токена  и возвращает объект токена

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
