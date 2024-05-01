package main

import (
	"api-gateway/transport"
	"api-gateway/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	cfg := utils.LoadConfig(".env")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		//-----| Public Routes
		// auth-service
		r.Post("/register", transport.ProxyRequest(cfg.AuthServiceURL, true, cfg))
		r.Post("/login", transport.ProxyRequest(cfg.AuthServiceURL, true, cfg))
		r.Post("/refreshToken", transport.ProxyRequest(cfg.AuthServiceURL, true, cfg))
		r.Post("/changePassword", transport.ProxyRequest(cfg.AuthServiceURL, true, cfg))

		// newsletter-management-service
		r.Get("/listNewsletters", transport.ProxyRequest(cfg.NewsletterServiceURL, true, cfg))

		// subscription-service
		r.Post("/subscribe", transport.ProxyRequest(cfg.SubscriptionServiceURL, true, cfg))
		r.Post("/unsubscribe", transport.ProxyRequest(cfg.SubscriptionServiceURL, true, cfg))

		//-----| Private Routes requiring JWT validation
		// newsletter-management-service
		r.Post("/createNewsletter", transport.ProxyRequest(cfg.NewsletterServiceURL, false, cfg))
		r.Post("/renameNewsletter", transport.ProxyRequest(cfg.NewsletterServiceURL, false, cfg))
		r.Post("/deleteNewsletter", transport.ProxyRequest(cfg.NewsletterServiceURL, false, cfg))

		// publishing-service
		r.Post("/publishPost", transport.ProxyRequest(cfg.PublishingServiceURL, false, cfg))
		r.Get("/listSubscribers", transport.ProxyRequest(cfg.PublishingServiceURL, false, cfg))
	})

	// Starting server
	address := cfg.IP + ":" + cfg.Port
	log.Printf("Server starting on %s", address)
	err := http.ListenAndServe(address, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
