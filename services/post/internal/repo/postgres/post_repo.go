package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/repoerrors"
	"github.com/sharlatan2005/posts_app_backend/services/post/internal/domain"
)

const tableName = "posts.posts"

type PostRepo struct {
	db *DB
}

func NewPostRepo(db *DB) *PostRepo {
	return &PostRepo{db: db}
}

func (r *PostRepo) Create(ctx context.Context, post *domain.Post) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (id, author_id, text)
		VALUES ($1, $2, $3)
		RETURNING created_at
	`, tableName)

	err := r.db.Pool.QueryRow(
		ctx,
		query,
		post.ID,
		post.AuthorID,
		post.Text,
	).Scan(&post.CreatedAt)

	if err != nil {
		return fmt.Errorf("Insert query: %w", repoerrors.ErrDB)
	}

	return nil
}

func (r *PostRepo) Delete(ctx context.Context, postID uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE id = $1
	`, tableName)

	result, err := r.db.Pool.Exec(
		ctx,
		query,
		postID)

	if err != nil {
		return fmt.Errorf("Delete user %s query: %w", postID.String(), err)
	}

	if rows := result.RowsAffected(); rows == 0 {
		return repoerrors.ErrNotFound
	}

	return nil
}

func (r *PostRepo) Update(ctx context.Context, postID uuid.UUID, newText string) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET text = $1
		WHERE id = $2
	`, tableName)

	result, err := r.db.Pool.Exec(
		ctx,
		query,
		newText,
		postID)

	if err != nil {
		return fmt.Errorf("Update user %s query: %w", postID.String(), err)
	}

	if rows := result.RowsAffected(); rows == 0 {
		return repoerrors.ErrNotFound
	}

	return nil
}
