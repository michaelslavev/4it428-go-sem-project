package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"newsletter-management-service/utils"
)

func main() {
	cfg := utils.LoadConfig(".env")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		// r.Post("/register", registerHandler)
	})

	// Starting server
	address := cfg.IP + ":" + cfg.Port
	log.Printf("Server starting on %s", address)
	err := http.ListenAndServe(address, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
