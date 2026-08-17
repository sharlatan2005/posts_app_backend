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

func (r *CommentRepo) GetAllPostComments(ctx context.Context, postID uuid.UUID) ([]*domain.Comment, error) {
	query := fmt.Sprintf(`
		SELECT id, post_id, author_id, text, created_at FROM %s
		WHERE post_id = $1
	`, tableName)

	rows, err := r.db.Pool.Query(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("Get all post %s comments: %w", postID.String(), err)
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		c := &domain.Comment{}
		err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.AuthorID,
			&c.Text,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("Scanning into comment: %w", err)
		}
		comments = append(comments, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *CommentRepo) DeleteAllByPost(ctx context.Context, postID uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE post_id = $1
	`, tableName)

	result, err := r.db.Pool.Exec(
		ctx,
		query,
		postID)

	if err != nil {
		return fmt.Errorf("Delete all comments from post %s query: %w", postID.String(), err)
	}

	if rows := result.RowsAffected(); rows == 0 {
		return repoerrors.ErrNotFound
	}

	return nil
}
