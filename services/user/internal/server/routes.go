package server

import (
	"net/http"

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
	r.mux.HandleFunc("POST /create_user", userHandler.CreateUser)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
