package main

import (
	"auth-service/utils"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

// User represents a user in the system
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Here, add logic to register the user in your system
	// This might include saving the user details in a database and handling passwords securely

	log.Printf("Registered user: %s", user.Username)
	w.WriteHeader(http.StatusOK)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Here, add logic to authenticate the user
	// This might include checking user credentials against a database

	log.Printf("Logged in user: %s", user.Username)
	w.WriteHeader(http.StatusOK)
}

func main() {
	cfg := utils.LoadConfig(".env")

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
