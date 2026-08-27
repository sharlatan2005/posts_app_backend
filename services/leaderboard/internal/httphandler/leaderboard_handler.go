package httphandler

import (
	"log"
	"net/http"

	"github.com/sharlatan2005/chat_app_go_backend_pkg/responseutils"
	"github.com/sharlatan2005/posts_app_backend/services/leaderboard/internal/service"
)

type LeaderboardHandler struct {
	leaderboardService *service.LeaderboardService
}

func NewLeaderboardHandler(leaderboardService *service.LeaderboardService) *LeaderboardHandler {
	return &LeaderboardHandler{
		leaderboardService: leaderboardService,
	}
}

func (h *LeaderboardHandler) GetLeaders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responseutils.StatusMethodNotAllowed(w, "Неправильный метод (должен быть GET)")
		return
	}

	leaders, err := h.leaderboardService.GetLeaders(r.Context())

	if err != nil {
		switch {
		default:
			log.Println(err.Error())
			responseutils.InternalServerError(w, "Внутренняя ошибка сервера")
			return
		}
	}

	responseutils.JSON(w, http.StatusCreated, leaders)
}
