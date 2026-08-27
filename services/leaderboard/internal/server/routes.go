package server

import (
	"net/http"

	"github.com/sharlatan2005/chat_app_go_backend_pkg/middleware"
	"github.com/sharlatan2005/posts_app_backend/services/leaderboard/internal/httphandler"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

func (r *Router) SetupRoutes(
	leaderboardHandler *httphandler.LeaderboardHandler,
) {
	r.mux.HandleFunc("GET /api/leaderboard/get_leaders",
		middleware.Chain(
			leaderboardHandler.GetLeaders,
			middleware.CORS,
			middleware.Logger))

}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
