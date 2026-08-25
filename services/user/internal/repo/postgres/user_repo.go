package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/repoerrors"
	"github.com/sharlatan2005/posts_app_backend/services/user/internal/domain"
)

const tableName = "users.users"

type UserRepo struct {
	db *DB
}

func NewUserRepo(db *DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (id, username, password_hash)
		VALUES ($1, $2, $3)
		ON CONFLICT (username) DO NOTHING
		RETURNING score, created_at
	`, tableName)

	result, err := r.db.Pool.Exec(ctx, query,
		user.ID,
		user.Username,
		user.Password_hash,
	)

	if err != nil {
		return fmt.Errorf("Insert query: %w", err)
	}

	if rows := result.RowsAffected(); rows == 0 {
		return repoerrors.ErrDuplicate
	}

	return nil
}

func (r *UserRepo) Exists(ctx context.Context, username string) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE username = $1)`, tableName)

	exists := false
	err := r.db.Pool.QueryRow(
		ctx,
		query,
		username,
	).Scan(
		&exists,
	)
	if err != nil {
		return false, fmt.Errorf("Query error: %w", err)
	}

	return exists, nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := fmt.Sprintf(
		`
		SELECT id, username, password_hash, name, surname, score, created_at
		FROM %s
		WHERE username = $1
	`, tableName)

	user := &domain.User{}

	err := r.db.Pool.QueryRow(
		ctx,
		query,
		username,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Password_hash,
		&user.Name,
		&user.Surname,
		&user.Score,
		&user.Created_at,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, repoerrors.ErrNotFound
		default:
			return nil, fmt.Errorf("Query of user: %w", err)
		}
	}

	return user, nil
}

func (r *UserRepo) GetByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	query := fmt.Sprintf(`
		SELECT username, password_hash, name, surname, score, created_at
		FROM %s
		WHERE id = $1
	`, tableName)

	user := &domain.User{
		ID: userID,
	}

	err := r.db.Pool.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&user.Username,
		&user.Password_hash,
		&user.Name,
		&user.Surname,
		&user.Score,
		&user.Created_at,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, repoerrors.ErrNotFound
		default:
			return nil, fmt.Errorf("Select user %s query: %w", userID.String(), err)
		}
	}

	return user, nil
}
