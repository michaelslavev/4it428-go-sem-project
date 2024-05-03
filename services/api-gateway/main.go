package main

import (
	"api-gateway/transport"
	"api-gateway/utils"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := utils.LoadConfig(".env")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		// auth-service
		r.Post("/register", transport.ProxyRequest(cfg.AuthServiceURL, true, cfg))
		r.Post("/login", transport.ProxyRequest(cfg.AuthServiceURL, true, cfg))
		r.Post("/refreshToken", transport.ProxyRequest(cfg.AuthServiceURL, true, cfg))
		
		// newsletter-management-service
		r.Post("/newsletters", transport.ProxyRequest(cfg.NewsletterServiceURL, false, cfg))
		r.Put("/newsletters", transport.ProxyRequest(cfg.NewsletterServiceURL, false, cfg))
		r.Get("/newsletters", transport.ProxyRequest(cfg.NewsletterServiceURL, true, cfg))
		r.Delete("/newsletters", transport.ProxyRequest(cfg.NewsletterServiceURL, false, cfg))
		r.Get("/subscribers", transport.ProxyRequest(cfg.PublishingServiceURL, false, cfg))

		// subscription-service
		r.Post("/subscribe", transport.ProxyRequest(cfg.SubscriptionServiceURL, false, cfg))
		r.Post("/unsubscribe", transport.ProxyRequest(cfg.SubscriptionServiceURL, false, cfg))

		// publishing-service
		r.Post("/posts", transport.ProxyRequest(cfg.PublishingServiceURL, false, cfg))
		r.Get("/posts", transport.ProxyRequest(cfg.PublishingServiceURL, false, cfg))
	})

	// Starting server
	address := cfg.IP + ":" + cfg.Port
	log.Printf("Server starting on %s", address)
	err := http.ListenAndServe(address, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
