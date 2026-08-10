package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/clients"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/clients/user"
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

type RegisterResult struct {
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *AuthService) RegisterUser(ctx context.Context, username string, password_hash string) (*RegisterResult, error) {
	err := s.userClient.CreateUser(ctx, username, password_hash)
	if err != nil {
		var extErr *clients.ExternalServiceError
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

	return &RegisterResult{
		Token:     token,
		UserID:    new_uuid.String(),
		Username:  username,
		CreatedAt: time.Now(),
	}, nil
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
