package main

import (
	"api-gateway/transport"
	"api-gateway/utils"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"strings"
)

func main() {
	cfg := utils.LoadConfig(".env")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Public Routes
	// auth-service
	r.Post("/register", transport.ProxyRequest(cfg.AuthServiceURL))
	r.Post("/login", transport.ProxyRequest(cfg.AuthServiceURL))
	r.Post("/refreshToken", transport.ProxyRequest(cfg.AuthServiceURL))
	r.Post("/changePassword", transport.ProxyRequest(cfg.AuthServiceURL))

	// newsletter-management-service
	r.Get("/listNewsletters", transport.ProxyRequest(cfg.NewsletterServiceURL))

	// subscription-service
	r.Post("/subscribe", transport.ProxyRequest(cfg.SubscriptionServiceURL))
	r.Post("/unsubscribe", transport.ProxyRequest(cfg.SubscriptionServiceURL))

	// Private Routes requiring JWT validation
	r.Route("/api", func(r chi.Router) {
		r.Use(JWTMiddleware)

		// newsletter-management-service
		r.Post("/createNewsletter", transport.ProxyRequest(cfg.NewsletterServiceURL))
		r.Post("/renameNewsletter", transport.ProxyRequest(cfg.NewsletterServiceURL))
		r.Post("/deleteNewsletter", transport.ProxyRequest(cfg.NewsletterServiceURL))

		// publishing-service
		r.Post("/publishPost", transport.ProxyRequest(cfg.PublishingServiceURL))
		r.Get("/listSubscribers", transport.ProxyRequest(cfg.PublishingServiceURL))
	})

	// Starting server
	address := cfg.IP + ":" + cfg.Port
	log.Printf("Server starting on %s", address)
	err := http.ListenAndServe(address, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

			return []byte("your-256-bit-secret"), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
