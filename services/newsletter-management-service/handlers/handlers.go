package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	supa "github.com/nedpals/supabase-go"
	"log"
	"net/http"
	"newsletter-management-service/handlers/model"
	"newsletter-management-service/handlers/sql"
	"newsletter-management-service/utils"
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
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
	w.WriteHeader(statusCode)
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

func (hd *CustomHandler) GetNewslettersHandler(w http.ResponseWriter, r *http.Request) {
	newsletters, err := hd.Repository.ListNewsletters(r.Context())
	if err != nil {
		handleError(w, "Failed to fetch newsletters", err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, newsletters, http.StatusOK)
}

func (hd *CustomHandler) CreateNewsletter(w http.ResponseWriter, r *http.Request) {
	token := utils.GetBearerToken(r)
	userUUId, _ := utils.ExtractSubFromToken(token)

	var newNewsletter model.NewNewsletter
	if !decodeRequest(w, r, &newNewsletter) {
		return
	}

	createdNewsletter, err := hd.Repository.CreateNewsletter(r.Context(), newNewsletter, userUUId)
	if err != nil {
		handleError(w, "Failed to create newsletters", err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, createdNewsletter, http.StatusOK)
}

func (hd *CustomHandler) RenameNewsletter(w http.ResponseWriter, r *http.Request) {
	token := utils.GetBearerToken(r)
	userUUId, _ := utils.ExtractSubFromToken(token)

	postId := chi.URLParam(r, "id")
	if postId == "" {
		handleError(w, "ID is required", nil, http.StatusBadRequest)
		return
	}

	var updatedNewsletter model.UpdateNewsletter
	if !decodeRequest(w, r, &updatedNewsletter) {
		return
	}

	updatedNewsletter.Id = postId

	uNewsletter, err := hd.Repository.RenameNewsletter(r.Context(), updatedNewsletter, userUUId)
	if err != nil {
		handleError(w, "Failed to update newsletter", err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, uNewsletter, http.StatusOK)
}

func (hd *CustomHandler) DeleteNewsletter(w http.ResponseWriter, r *http.Request) {
	token := utils.GetBearerToken(r)
	userUUId, _ := utils.ExtractSubFromToken(token)

	postId := chi.URLParam(r, "id")
	if postId == "" {
		handleError(w, "ID is required", nil, http.StatusBadRequest)
		return
	}

	err := hd.Repository.DeleteNewsletter(r.Context(), postId, userUUId)
	if err != nil {
		handleError(w, "Failed to delete newsletter", err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, "", http.StatusNoContent)
}

func (hd *CustomHandler) GetNewsletterSubscribers(w http.ResponseWriter, r *http.Request) {
	token := utils.GetBearerToken(r)
	userUUId, _ := utils.ExtractSubFromToken(token)

	postId := chi.URLParam(r, "id")
	if postId == "" {
		handleError(w, "ID is required", nil, http.StatusBadRequest)
		return
	}

	subscribers, err := hd.Repository.GetNewsletterSubscribers(r.Context(), postId, userUUId)
	if err != nil {
		handleError(w, "Failed to fetch subscribers", err, http.StatusInternalServerError)
		return
	}
	log.Printf("Subscribers: %v", subscribers)
	sendJSON(w, subscribers, http.StatusOK)
}
