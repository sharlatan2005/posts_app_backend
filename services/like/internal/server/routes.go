package server

import (
	"net/http"

	"github.com/sharlatan2005/chat_app_go_backend_pkg/auth"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/middleware"
	"github.com/sharlatan2005/posts_app_backend/services/like/internal/httphandler"
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
	likeHandler *httphandler.LikeHandler,
) {
	r.mux.HandleFunc("POST /api/like/add_like",
		middleware.Chain(
			likeHandler.AddLike,
			authStruct.Authenticate,
			middleware.CORS,
			middleware.Logger))

	r.mux.HandleFunc("DELETE /api/like/remove_like",
		middleware.Chain(
			likeHandler.RemoveLike,
			authStruct.Authenticate,
			middleware.CORS,
			middleware.Logger))

	r.mux.HandleFunc("GET /api/like/get_liked_user_ids",
		middleware.Chain(
			likeHandler.GetLikedUserIDs,
			authStruct.Authenticate,
			middleware.CORS,
			middleware.Logger))

	r.mux.HandleFunc("DELETE /delete_likes_by_post",
		middleware.Chain(
			likeHandler.DeleteLikesByPost,
			middleware.CORS,
			middleware.Logger))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
