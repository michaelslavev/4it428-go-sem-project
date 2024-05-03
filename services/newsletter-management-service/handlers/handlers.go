package handlers

import (
	"encoding/json"
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

	err := hd.Repository.CreateNewsletter(r.Context(), newNewsletter, userUUId)
	if err != nil {
		handleError(w, "Failed to create newsletters", err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, "{}", http.StatusOK)
}

func (hd *CustomHandler) RenameNewsletter(w http.ResponseWriter, r *http.Request) {
	// parse url to get id
	// parse request body to get new name
	// call repository to update newsletter name
	// return response
	panic("not implemented")
}

func (hd *CustomHandler) DeleteNewsletter(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}
