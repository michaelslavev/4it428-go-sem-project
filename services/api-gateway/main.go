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
		r.Post("/register", transport.ProxyRequest(cfg.AuthServiceURL, true))
		r.Post("/login", transport.ProxyRequest(cfg.AuthServiceURL, true))
		r.Post("/refreshToken", transport.ProxyRequest(cfg.AuthServiceURL, true))
		r.Post("/changePassword", transport.ProxyRequest(cfg.AuthServiceURL, true))

		// newsletter-management-service
		r.Get("/listNewsletters", transport.ProxyRequest(cfg.NewsletterServiceURL, true))

		// subscription-service
		r.Post("/subscribe", transport.ProxyRequest(cfg.SubscriptionServiceURL, true))
		r.Post("/unsubscribe", transport.ProxyRequest(cfg.SubscriptionServiceURL, true))

		//-----| Private Routes requiring JWT validation
		// newsletter-management-service
		r.Post("/createNewsletter", transport.ProxyRequest(cfg.NewsletterServiceURL, false))
		r.Post("/renameNewsletter", transport.ProxyRequest(cfg.NewsletterServiceURL, false))
		r.Post("/deleteNewsletter", transport.ProxyRequest(cfg.NewsletterServiceURL, false))

		// publishing-service
		r.Post("/publishPost", transport.ProxyRequest(cfg.PublishingServiceURL, false))
		r.Get("/listSubscribers", transport.ProxyRequest(cfg.PublishingServiceURL, false))
	})

	// Starting server
	address := cfg.IP + ":" + cfg.Port
	log.Printf("Server starting on %s", address)
	err := http.ListenAndServe(address, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
