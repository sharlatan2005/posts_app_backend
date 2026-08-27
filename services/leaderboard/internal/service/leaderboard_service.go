package service

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/sharlatan2005/posts_app_backend/services/leaderboard/internal/cache"
	"github.com/sharlatan2005/posts_app_backend/services/leaderboard/internal/domain"
	"github.com/sharlatan2005/posts_app_backend/services/leaderboard/internal/repo"
)

type LeaderboardService struct {
	leaderboardRepo  repo.LeaderboardRepo
	leaderboardCache cache.LeaderboardCache
}

func NewLeaderboardService(
	leaderboardRepo repo.LeaderboardRepo,
	leaderboardCache cache.LeaderboardCache) *LeaderboardService {
	return &LeaderboardService{
		leaderboardRepo:  leaderboardRepo,
		leaderboardCache: leaderboardCache,
	}
}

func (s *LeaderboardService) UpgradeScore(
	ctx context.Context,
	userID uuid.UUID,
	points int,
) error {
	if err := s.leaderboardRepo.UpgradeScore(ctx, userID, points); err != nil {
		return err
	}

	if err := s.leaderboardCache.UpgradeScore(ctx, userID, points); err != nil {
		log.Printf("Redis score update failed for user %s: %v", userID, err)
	}

	return nil
}

func (s *LeaderboardService) GetLeaders(ctx context.Context) ([]domain.LeaderboardItem, error) {
	leaders, err := s.leaderboardCache.GetLeaders(ctx)
	if err != nil {
		switch {
		default:
			return nil, fmt.Errorf("Getting leaders: %w", err)
		}
	}

	return leaders, nil
}
