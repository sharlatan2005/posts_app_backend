package cache

import (
	"context"

	"github.com/google/uuid"
	"github.com/sharlatan2005/posts_app_backend/services/leaderboard/internal/domain"
)

type LeaderboardCache interface {
	UpgradeScore(ctx context.Context, userID uuid.UUID, points int) error
	GetLeaders(ctx context.Context) ([]domain.LeaderboardItem, error)
}
