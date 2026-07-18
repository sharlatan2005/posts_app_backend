package postgres

import (
	"context"
	"fmt"

	"github.com/sharlatan2005/posts_app_backend/pkg/repoerrors"
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
		INSERT INTO %s (id, username, password_hash, name, surname)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (username) DO NOTHING
		RETURNING score, created_at
	`, tableName)

	result, err := r.db.Pool.Exec(ctx, query,
		user.ID,
		user.Username,
		user.Password_hash,
		user.Name,
		user.Surname,
	)

	if err != nil {
		return fmt.Errorf("Insert query: %w", repoerrors.ErrDB)
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
