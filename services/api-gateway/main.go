package main

import (
	"log"
	"net/http"

	"github.com/sharlatan2005/posts_app_backend/services/api-gateway/internal/config"
	"github.com/sharlatan2005/posts_app_backend/services/api-gateway/internal/proxy"
)

func main() {
	cfg := config.Load()

	authProxy := proxy.MustNewProxy(cfg.AuthServiceHost, cfg.AuthServicePort)
	userProxy := proxy.MustNewProxy(cfg.UserServiceHost, cfg.UserServicePort)

	// API
	http.HandleFunc("api/auth/", proxy.ProxyHandler(authProxy))
	http.HandleFunc("api/user/", proxy.ProxyHandler(userProxy))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
