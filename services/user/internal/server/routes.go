package server

import (
	"net/http"

	"github.com/sharlatan2005/chat_app_go_backend_pkg/middleware"
	"github.com/sharlatan2005/posts_app_backend/services/user/internal/httphandler"
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
	userHandler *httphandler.UserHandler,
) {
	r.mux.HandleFunc("POST /create_user",
		middleware.Chain(
			userHandler.CreateUser,
			middleware.CORS,
			middleware.Logger))

	r.mux.HandleFunc("GET /get_user_by_username",
		middleware.Chain(
			userHandler.GetByUsername,
			middleware.CORS,
			middleware.Logger))

	r.mux.HandleFunc("GET /api/user/get_user_by_id",
		middleware.Chain(
			userHandler.GetUserByID,
			middleware.CORS,
			middleware.Logger))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
