package transport

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func ProxyRequest(target string, isPublic bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check JWT if the route is not public
		if !isPublic {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				http.Error(w, "Authorization token required", http.StatusUnauthorized)
				return
			}

			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				// Ensure the secret used here is the same one used when issuing the JWT
				return []byte("your-256-bit-secret"), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
		}

		// Proxy REQ to service
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
