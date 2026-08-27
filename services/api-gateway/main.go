package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sharlatan2005/posts_app_backend/services/api-gateway/internal/config"
	"github.com/sharlatan2005/posts_app_backend/services/api-gateway/internal/proxy"
)

func main() {
	cfg := config.Load()

	authProxy := proxy.MustNewProxy(cfg.AuthServiceHost, cfg.AuthServicePort)
	userProxy := proxy.MustNewProxy(cfg.UserServiceHost, cfg.UserServicePort)
	postProxy := proxy.MustNewProxy(cfg.PostServiceHost, cfg.PostServicePort)
	commentProxy := proxy.MustNewProxy(cfg.CommentServiceHost, cfg.CommentServicePort)
	likeProxy := proxy.MustNewProxy(cfg.LikeServiceHost, cfg.LikeServicePort)
	leaderboardProxy := proxy.MustNewProxy(cfg.LeaderboardServiceHost, cfg.LeaderboardServicePort)

	// API
	http.HandleFunc("/api/auth/", proxy.ProxyHandler(authProxy))
	http.HandleFunc("/api/user/", proxy.ProxyHandler(userProxy))
	http.HandleFunc("/api/post/", proxy.ProxyHandler(postProxy))
	http.HandleFunc("/api/comment/", proxy.ProxyHandler(commentProxy))
	http.HandleFunc("/api/like/", proxy.ProxyHandler(likeProxy))
	http.HandleFunc("/api/leaderboard/", proxy.ProxyHandler(leaderboardProxy))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.ServicePort), nil))
}
