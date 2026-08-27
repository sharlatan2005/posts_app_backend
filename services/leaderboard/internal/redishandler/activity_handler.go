package redishandler

import (
	"context"
	"encoding/json"

	"github.com/sharlatan2005/chat_app_go_backend_pkg/events"
	"github.com/sharlatan2005/posts_app_backend/services/leaderboard/internal/service"
)

type ActivityHandler struct {
	service *service.LeaderboardService
}

func NewActivityHandler(service *service.LeaderboardService) *ActivityHandler {
	return &ActivityHandler{
		service: service,
	}
}

func (h *ActivityHandler) Handle(msg []byte) error {
	var activity events.Activity
	if err := json.Unmarshal(msg, &activity); err != nil {
		return err
	}

	points := map[string]int{
		"like":    1,
		"comment": 5,
		"post":    10,
	}[activity.Type]

	if points == 0 {
		return nil
	}

	return h.service.UpgradeScore(
		context.Background(),
		activity.UserID,
		points,
	)
}
