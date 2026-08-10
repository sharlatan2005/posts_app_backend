package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func MustNewProxy(host, port string) *httputil.ReverseProxy {
	u := &url.URL{Scheme: "http", Host: host + ":" + port}
	return httputil.NewSingleHostReverseProxy(u)
}

func ProxyHandler(p *httputil.ReverseProxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p.ServeHTTP(w, r)
	}
}
