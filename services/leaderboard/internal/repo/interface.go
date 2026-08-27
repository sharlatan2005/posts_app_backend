package repo

import (
	"context"

	"github.com/google/uuid"
)

type LeaderboardRepo interface {
	UpgradeScore(ctx context.Context, userID uuid.UUID, points int) error
}
