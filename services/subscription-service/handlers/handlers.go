package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"subscription-service/handlers/sql"

	supa "github.com/nedpals/supabase-go"
)

type CustomHandler struct {
	SupabaseClient *supa.Client
	Repository     *sql.Repository
}

func NewCustomHandler(client *supa.Client, repository *sql.Repository) *CustomHandler {
	return &CustomHandler{
		SupabaseClient: client,
		Repository:     repository,
	}
}

func sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func handleError(w http.ResponseWriter, message string, err error, statusCode int) {
	log.Printf("%s: %v", message, err)
	http.Error(w, message, statusCode)
}

func decodeRequest(w http.ResponseWriter, r *http.Request, dest interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		handleError(w, err.Error(), err, http.StatusBadRequest)
		return false
	}
	return true
}

func (hd *CustomHandler) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	subscriptions, err := hd.Repository.ListSubscriptions(r.Context())
	if err != nil {
		handleError(w, "Failed to fetch subscriptions", err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, subscriptions, http.StatusOK)
}