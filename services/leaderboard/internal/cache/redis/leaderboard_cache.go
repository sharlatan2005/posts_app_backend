package redis

import (
	"context"

	"github.com/google/uuid"
	"github.com/sharlatan2005/posts_app_backend/services/leaderboard/internal/domain"
)

type LeaderboardCache struct {
	db *DB
}

func NewLeaderboardCache(db *DB) *LeaderboardCache {
	return &LeaderboardCache{db: db}
}

const leaderboardKey = "leaderboard:scores"

func (c *LeaderboardCache) UpgradeScore(ctx context.Context, userID uuid.UUID, points int) error {

	// Атомарно добавляем очки в sorted set
	return c.db.Client.ZIncrBy(ctx, leaderboardKey, float64(points), userID.String()).Err()
}

func (c *LeaderboardCache) GetLeaders(ctx context.Context) ([]domain.LeaderboardItem, error) {
	items, err := c.db.Client.
		ZRevRangeWithScores(ctx, leaderboardKey, 0, 9).
		Result()

	if err != nil {
		return nil, err
	}

	result := make([]domain.LeaderboardItem, 0, len(items))

	for _, item := range items {
		result = append(result, domain.LeaderboardItem{
			UserID: item.Member.(string),
			Score:  int64(item.Score),
		})
	}

	return result, nil

}
