package server

import (
	"net/http"

	"github.com/sharlatan2005/posts_app_backend/services/auth/internal/httphandler"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

func (r *Router) SetupRoutes(authHandler *httphandler.AuthHandler) {
	r.mux.HandleFunc("POST /register", authHandler.Register)
}
