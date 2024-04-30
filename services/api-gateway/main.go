package main

import (
	"api-gateway/utils"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	cfg := utils.LoadConfig(".env")

	r := mux.NewRouter()
	address := cfg.IP + ":" + cfg.Port
	log.Printf("Server starting on %s", address)

	// ----| PUBLIC ROUTES
	// auth-service
	r.HandleFunc("/register", proxyRequest(cfg.AuthServiceURL)).Methods("POST")
	r.HandleFunc("/login", proxyRequest(cfg.AuthServiceURL)).Methods("POST")
	r.HandleFunc("/refreshToken", proxyRequest(cfg.AuthServiceURL)).Methods("POST")
	r.HandleFunc("/changePassword", proxyRequest(cfg.AuthServiceURL)).Methods("POST")

	// newsletter-management-service
	r.HandleFunc("/listNewsletters", proxyRequest(cfg.NewsletterServiceURL)).Methods("GET")

	// subscription-service
	r.HandleFunc("/subscribe", proxyRequest(cfg.SubscriptionServiceURL)).Methods("POST")
	r.HandleFunc("/unsubscribe", proxyRequest(cfg.SubscriptionServiceURL)).Methods("POST")

	// ----| PRIVATE ROUTES requiring JWT validation
	s := r.PathPrefix("/api").Subrouter()
	s.Use(JWTMiddleware)
	// newsletter-management-service
	s.HandleFunc("/createNewsletter", proxyRequest(cfg.NewsletterServiceURL)).Methods("POST")
	s.HandleFunc("/renameNewsletter", proxyRequest(cfg.NewsletterServiceURL)).Methods("POST")
	s.HandleFunc("/deleteNewsletter", proxyRequest(cfg.NewsletterServiceURL)).Methods("POST")

	// publishing-service
	s.HandleFunc("/publishPost", proxyRequest(cfg.PublishingServiceURL)).Methods("POST")
	s.HandleFunc("/listSubscribers", proxyRequest(cfg.PublishingServiceURL)).Methods("GET")

	// Starting server
	err := http.ListenAndServe(address, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func proxyRequest(target string) func(http.ResponseWriter, *http.Request) {
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

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		// Trim Bearer prefix if you use it
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte("your-256-bit-secret"), nil // Use your secret
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Token is valid
		next.ServeHTTP(w, r)
	})
}
