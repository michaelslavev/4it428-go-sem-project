package handlers

import (
	"auth-service/handlers/model"
	"context"
	"encoding/json"
	supa "github.com/nedpals/supabase-go"
	"log"
	"net/http"
	"strings"
)

type CustomHandler struct {
	SupabaseClient *supa.Client
}

func NewCustomHandler(client *supa.Client) *CustomHandler {
	return &CustomHandler{
		SupabaseClient: client,
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

func (hd *CustomHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if !decodeRequest(w, r, &user) {
		return
	}

	resp, err := hd.SupabaseClient.Auth.SignUp(r.Context(), supa.UserCredentials{
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		handleError(w, "Failed to register user", err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, resp, http.StatusOK)
}

func (hd *CustomHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if !decodeRequest(w, r, &user) {
		return
	}

	loggedUser, err := hd.SupabaseClient.Auth.SignIn(r.Context(), supa.UserCredentials{
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		handleError(w, "Failed to login user", err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, loggedUser, http.StatusOK)
}

func (hd *CustomHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var user model.LoggedUserInfo
	if !decodeRequest(w, r, &user) {
		return
	}

	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	refreshedUser, err := hd.SupabaseClient.Auth.RefreshUser(r.Context(), tokenString, user.RefreshToken)
	if err != nil {
		handleError(w, "Failed to refresh token for user", err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, refreshedUser, http.StatusOK)
}

func (hd *CustomHandler) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if !decodeRequest(w, r, &user) {
		return
	}

	ctx := context.Background()
	err := hd.SupabaseClient.Auth.ResetPasswordForEmail(ctx, user.Email)
	if err != nil {
		handleError(w, "Failed to reset password for user", err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]string{"message": "Password reset email sent"}, http.StatusOK)
}
