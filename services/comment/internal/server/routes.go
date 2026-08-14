package server

import (
	"net/http"

	"github.com/sharlatan2005/chat_app_go_backend_pkg/auth"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/middleware"
	"github.com/sharlatan2005/posts_app_backend/services/comment/internal/httphandler"
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
	authStruct *auth.Auth,
	commentHandler *httphandler.CommentHandler,
) {
	r.mux.HandleFunc("POST /api/comment/create_comment",
		middleware.Chain(
			commentHandler.CreateComment,
			authStruct.Authenticate,
			middleware.CORS,
			middleware.Logger))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
