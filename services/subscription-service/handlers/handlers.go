package handlers

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"subscription-service/handlers/sql"
	"subscription-service/utils"

	"github.com/go-chi/chi/v5"
	supa "github.com/nedpals/supabase-go"
	"github.com/resend/resend-go/v2"
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

func loadEmailTemplate(data map[string]string) (string, error) {
    fileData, err := fs.ReadFile(SubscribedMailTemplate, "static/mail_template.html")
    if err != nil {
        return "", err
    }

    var builder strings.Builder
    scanner := bufio.NewScanner(strings.NewReader(string(fileData)))
    for scanner.Scan() {
        line := scanner.Text()
        for key, value := range data {
            placeholder := fmt.Sprintf("{{%s}}", key)
            line = strings.ReplaceAll(line, placeholder, value)
        }
        builder.WriteString(line + "\n")
    }

    if err := scanner.Err(); err != nil {
        return "", err
    }

    return builder.String(), nil
}

func (hd *CustomHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	token := utils.GetBearerToken(r)
	userId, _ := utils.ExtractSubFromToken(token)

	newsletterId := chi.URLParam(r, "id")
	if newsletterId == "" {
		handleError(w, "ID is required", nil, http.StatusBadRequest)
		return
	}

	subscription, subscriber, newsletter, err := hd.Repository.Subscribe(r.Context(), newsletterId, userId)
	if err != nil {
        if errors.Is(err, utils.ErrSubscriptionExists) {
            handleError(w, err.Error(), nil, http.StatusConflict) // 409 Conflict
        } else {
            handleError(w, "failed to subscribe", err, http.StatusInternalServerError)
        }
        return
    }

	cfg := utils.LoadConfig(".env")
	apiKey := cfg.ResendApiKey
	client := resend.NewClient(apiKey)
	// TODO: Change to production url
	unsubscribeUrl := "http://localhost:9069/api/unsubscribe"

    data := map[string]string{
        "unsubscribeUrl": unsubscribeUrl,
        "newsletterId":   newsletterId,
        "userId":         userId,
    }

    mail, loadErr := loadEmailTemplate(data)

    if loadErr != nil {
		log.Printf("Failed load mail template: %v", err)
    }

	params := &resend.SendEmailRequest{
		From:    "newsletter@tapeer.cz",
		To:      []string{subscriber.Email},
		Subject: fmt.Sprintf("You have been subscribed to %s.", newsletter.Title),
		Html: mail,
		Headers: map[string]string{
			"List-Unsubscribe": fmt.Sprintf("<%s?newsletterId=%s&userId=%s>", unsubscribeUrl, newsletterId, userId),
			"List-Unsubscribe-Post": "List-Unsubscribe=One-Click",
		},
	}

	_, emailErr := client.Emails.Send(params)
	if emailErr != nil {
		log.Printf("Failed to send email to %s: %v", subscriber.Email, err)
	}

	sendJSON(w, subscription, http.StatusOK)
}

func (hd *CustomHandler) Unsubcribe(w http.ResponseWriter, r *http.Request) {
	newsletterId := r.URL.Query().Get("newsletterId")
    userId := r.URL.Query().Get("userId")

	if newsletterId == "" || userId == "" {
        http.Error(w, "Missing newsletterId or userId", http.StatusBadRequest)
        return
    }

	err := hd.Repository.Unsubcribe(r.Context(), newsletterId, userId)
	if err != nil {
		handleError(w, "failed to unsubcribe", err, http.StatusInternalServerError)
		return
	}

	unsubbedPage, readErr := fs.ReadFile(Unsubscribed, "static/unsubscribed.html")
    if readErr != nil {
        handleError(w, "Failed to load the unsubscribe page", readErr, http.StatusInternalServerError)
        return
    }

	w.Header().Set("Content-Type", "text/html")
	w.Write(unsubbedPage)
}