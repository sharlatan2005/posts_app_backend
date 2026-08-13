package server

import (
	"net/http"

	"github.com/sharlatan2005/chat_app_go_backend_pkg/auth"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/middleware"
	"github.com/sharlatan2005/posts_app_backend/services/post/internal/httphandler"
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
	postHandler *httphandler.PostHandler,
) {
	r.mux.HandleFunc("POST /create_post",
		middleware.Chain(
			postHandler.CreatePost,
			authStruct.Authenticate,
			middleware.CORS,
			middleware.Logger))

	r.mux.HandleFunc("DELETE /delete_post",
		middleware.Chain(
			postHandler.DeletePost,
			authStruct.Authenticate,
			middleware.CORS,
			middleware.Logger))

	r.mux.HandleFunc("PUT /update_post",
		middleware.Chain(
			postHandler.UpdatePost,
			authStruct.Authenticate,
			middleware.CORS,
			middleware.Logger))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
