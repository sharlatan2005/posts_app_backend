package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/repoerrors"
)

const tableName = "users.users"

type LeaderboardRepo struct {
	db *DB
}

func NewLeaderboardRepo(db *DB) *LeaderboardRepo {
	return &LeaderboardRepo{db: db}
}

func (r *LeaderboardRepo) UpgradeScore(ctx context.Context, userID uuid.UUID, points int) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET score = score + $1
		WHERE id = $2
	`, tableName)

	result, err := r.db.Pool.Exec(
		ctx,
		query,
		points,
		userID,
	)

	if err != nil {
		return fmt.Errorf("upgrade score of user %s: %w", userID.String(), err)
	}

	if result.RowsAffected() == 0 {
		return repoerrors.ErrNotFound
	}

	return nil
}
