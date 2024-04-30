package transport

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func ProxyRequest(target string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		targetUrl, err := url.Parse(target)
		if err != nil {
			log.Printf("Error parsing URL: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(targetUrl)
		r.URL.Host = targetUrl.Host
		r.URL.Scheme = targetUrl.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = targetUrl.Host

		proxy.ServeHTTP(w, r)
	}
}
