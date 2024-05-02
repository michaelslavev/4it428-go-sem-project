package main

import (
	"auth-service/handlers"
	"auth-service/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	supa "github.com/nedpals/supabase-go"
	"log"
	"net/http"
)

var supabase *supa.Client
var cfg utils.ServerConfig

func init() {
	cfg = utils.LoadConfig(".env")
	supabase = supa.CreateClient(cfg.SupabaseURL, cfg.SupabaseKEY)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	hd := handlers.NewCustomHandler(supabase)

	r.Route("/api", func(r chi.Router) {
		r.Post("/register", hd.RegisterHandler)
		r.Post("/login", hd.LoginHandler)
		r.Post("/refreshToken", hd.RefreshTokenHandler)

		// Sad, we would need UI for that
		//r.Post("/resetPassword", hd.ResetPasswordHandler)
	})

	// Starting server
	address := cfg.IP + ":" + cfg.Port
	log.Printf("Server starting on %s", address)
	err := http.ListenAndServe(address, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
