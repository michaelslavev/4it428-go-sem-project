package main

import (
	"auth-service/utils"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	supa "github.com/nedpals/supabase-go"
	"log"
	"net/http"
)

// User represents a user in the system
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var supabase *supa.Client
var cfg utils.ServerConfig

func init() {
	cfg = utils.LoadConfig(".env")
	supabase = supa.CreateClient(cfg.SupabaseURL, cfg.SupabaseKEY)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	resp, err := supabase.Auth.SignUp(ctx, supa.UserCredentials{
		Email:    user.Username,
		Password: user.Password,
	})
	if err != nil {
		log.Printf("Failed to register user: %v", err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	log.Printf("Registered user: %s", user.Username)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	loggedUser, err := supabase.Auth.SignIn(ctx, supa.UserCredentials{
		Email:    user.Username,
		Password: user.Password,
	})
	if err != nil {
		log.Printf("Failed to login user: %v", err)
		http.Error(w, "Failed to login user", http.StatusUnauthorized)
		return
	}

	log.Printf("Registered user: %s", user.Username)
	err = json.NewEncoder(w).Encode(loggedUser)
	if err != nil {
		return
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Post("/register", registerHandler)
		r.Post("/login", loginHandler)
		r.Post("/refreshToken", loginHandler)
		r.Post("/changePassword", loginHandler)
	})

	// Starting server
	address := cfg.IP + ":" + cfg.Port
	log.Printf("Server starting on %s", address)
	err := http.ListenAndServe(address, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
