package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"subscription-service/handlers/sql"
	"subscription-service/utils"

	"github.com/go-chi/chi/v5"
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

func (hd *CustomHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	token := utils.GetBearerToken(r)
	userUUId, _ := utils.ExtractSubFromToken(token)

	newsletterId := chi.URLParam(r, "id")
	if newsletterId == "" {
		handleError(w, "ID is required", nil, http.StatusBadRequest)
		return
	}

	subscription, err := hd.Repository.Subscribe(r.Context(), newsletterId, userUUId)
	if err != nil {
        if errors.Is(err, utils.ErrSubscriptionExists) {
            handleError(w, err.Error(), nil, http.StatusConflict) // 409 Conflict
        } else {
            handleError(w, "failed to subscribe", err, http.StatusInternalServerError)
        }
        return
    }

	sendJSON(w, subscription, http.StatusOK)
}

func (hd *CustomHandler) Unsubcribe(w http.ResponseWriter, r *http.Request) {
	token := utils.GetBearerToken(r)
	userUUId, _ := utils.ExtractSubFromToken(token)

	newsletterId := chi.URLParam(r, "id")
	if newsletterId == "" {
		handleError(w, "ID is required", nil, http.StatusBadRequest)
		return
	}

	err := hd.Repository.Unsubcribe(r.Context(), newsletterId, userUUId)
	if err != nil {
		handleError(w, "Failed to unsubcribe", err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, "", http.StatusNoContent)
}