package main

import (
	"api-gateway/transport"
	"api-gateway/utils"
	_ "embed"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

//go:embed api/openapi.yaml
var OpenAPI string

func main() {
	cfg := utils.LoadConfig(".env")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		// openapi specification file
		r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/yaml")
			if _, err := w.Write([]byte(OpenAPI)); err != nil {
				http.Error(w, "Error sending response", http.StatusInternalServerError)
				return
			}
		})

		// auth-service
		r.Post("/register", transport.ProxyRequest(cfg.AuthServiceURL, true, cfg))
		r.Post("/login", transport.ProxyRequest(cfg.AuthServiceURL, true, cfg))
		r.Post("/refreshToken", transport.ProxyRequest(cfg.AuthServiceURL, true, cfg))

		// newsletter-management-service
		r.Post("/newsletters", transport.ProxyRequest(cfg.NewsletterServiceURL, false, cfg))
		r.Put("/newsletters/{id}", transport.ProxyRequest(cfg.NewsletterServiceURL, false, cfg))
		r.Get("/newsletters", transport.ProxyRequest(cfg.NewsletterServiceURL, true, cfg))
		r.Delete("/newsletters/{id}", transport.ProxyRequest(cfg.NewsletterServiceURL, false, cfg))
		r.Get("/subscribers/{id}", transport.ProxyRequest(cfg.NewsletterServiceURL, false, cfg))

		// subscription-service
		r.Post("/subscribe/{id}", transport.ProxyRequest(cfg.SubscriptionServiceURL, false, cfg))
		r.Get("/unsubscribe/{magic}", transport.ProxyRequest(cfg.SubscriptionServiceURL, true, cfg))

		// publishing-service
		r.Post("/posts", transport.ProxyRequest(cfg.PublishingServiceURL, false, cfg))
		r.Get("/posts", transport.ProxyRequest(cfg.PublishingServiceURL, false, cfg))
	})

	r.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// Starting server
	address := cfg.IP + ":" + cfg.Port
	log.Printf("Server starting on %s", address)
	err := http.ListenAndServe(address, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
