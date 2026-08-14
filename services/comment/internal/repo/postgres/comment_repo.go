package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/repoerrors"
	"github.com/sharlatan2005/posts_app_backend/services/comment/internal/domain"
)

const tableName = "comments.comments"

type CommentRepo struct {
	db *DB
}

func NewCommentRepo(db *DB) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) Create(ctx context.Context, comment *domain.Comment) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (id, post_id, author_id, text)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`, tableName)

	err := r.db.Pool.QueryRow(
		ctx,
		query,
		comment.ID,
		comment.PostID,
		comment.AuthorID,
		comment.Text,
	).Scan(&comment.CreatedAt)

	if err != nil {
		return fmt.Errorf("Insert query: %w", repoerrors.ErrDB)
	}

	return nil
}

func (r *CommentRepo) GetByID(ctx context.Context, commentID uuid.UUID) (*domain.Comment, error) {
	query := fmt.Sprintf(`
		SELECT post_id, author_id, text, created_at
		FROM %s
		WHERE id = $1
	`, tableName)

	comment := &domain.Comment{
		ID: commentID,
	}

	err := r.db.Pool.QueryRow(
		ctx,
		query,
		commentID,
	).Scan(
		&comment.PostID,
		&comment.AuthorID,
		&comment.Text,
		&comment.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, repoerrors.ErrNotFound
		default:
			return nil, fmt.Errorf("Select comment %s query: %w", commentID.String(), err)
		}
	}

	return comment, nil
}

func (r *CommentRepo) Delete(ctx context.Context, commentID uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE id = $1
	`, tableName)

	result, err := r.db.Pool.Exec(
		ctx,
		query,
		commentID)

	if err != nil {
		return fmt.Errorf("Delete comment %s query: %w", commentID.String(), err)
	}

	if rows := result.RowsAffected(); rows == 0 {
		return repoerrors.ErrNotFound
	}

	return nil
}

func (r *CommentRepo) Update(ctx context.Context, commentID uuid.UUID, newText string) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET text = $1
		WHERE id = $2
	`, tableName)

	result, err := r.db.Pool.Exec(
		ctx,
		query,
		newText,
		commentID)

	if err != nil {
		return fmt.Errorf("Update comment %s query: %w", commentID.String(), err)
	}

	if rows := result.RowsAffected(); rows == 0 {
		return repoerrors.ErrNotFound
	}

	return nil
}
