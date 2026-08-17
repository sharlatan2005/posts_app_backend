package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/repoerrors"
	"github.com/sharlatan2005/posts_app_backend/services/like/internal/domain"
)

const tableName = "likes.likes"

type LikeRepo struct {
	db *DB
}

func NewLikeRepo(db *DB) *LikeRepo {
	return &LikeRepo{db: db}
}

func (r *LikeRepo) AddLike(ctx context.Context, like *domain.Like) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (id, post_id, liker_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (post_id, liker_id) DO NOTHING
		RETURNING liked_at
	`, tableName)

	err := r.db.Pool.QueryRow(
		ctx,
		query,
		like.ID,
		like.PostID,
		like.LikerID,
	).Scan(&like.LikedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("Insert query: %w", err)
	}

	return nil
}

func (r *LikeRepo) RemoveLike(ctx context.Context, userID uuid.UUID, postID uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE liker_id = $1 and post_id = $2 
	`, tableName)

	result, err := r.db.Pool.Exec(
		ctx,
		query,
		userID,
		postID)

	if err != nil {
		return fmt.Errorf("Delete like from post %s query: %w", postID.String(), err)
	}

	if rows := result.RowsAffected(); rows == 0 {
		return repoerrors.ErrNotFound
	}

	return nil
}

func (r *LikeRepo) GetLikedUserIDs(ctx context.Context, postID uuid.UUID) ([]uuid.UUID, error) {
	query := fmt.Sprintf(`
		SELECT liker_id FROM %s
		WHERE post_id = $1
	`, tableName)

	rows, err := r.db.Pool.Query(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("Get all users that liked post %s: %w", postID.String(), err)
	}
	defer rows.Close()

	var likedUserIDs []uuid.UUID
	for rows.Next() {
		like := &domain.Like{}
		err := rows.Scan(
			&like.LikerID,
		)
		if err != nil {
			return nil, fmt.Errorf("Scanning into comment: %w", err)
		}
		likedUserIDs = append(likedUserIDs, like.LikerID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return likedUserIDs, nil
}

func (r *LikeRepo) DeleteAllByPost(ctx context.Context, postID uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE post_id = $1
	`, tableName)

	result, err := r.db.Pool.Exec(
		ctx,
		query,
		postID)

	if err != nil {
		return fmt.Errorf("Delete all likes from post %s query: %w", postID.String(), err)
	}

	if rows := result.RowsAffected(); rows == 0 {
		return repoerrors.ErrNotFound
	}

	return nil
}
