package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"publishing-service/handlers/model"
	"publishing-service/handlers/sql"
	"publishing-service/utils"
	"strings"

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

func decodeRequest(w http.ResponseWriter, r *http.Request, dest interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		handleError(w, err.Error(), err, http.StatusBadRequest)
		return false
	}
	return true
}

func loadEmailTemplate(data map[string]string) (string, error) {
    fileData, err := fs.ReadFile(PostTemplate, "static/post_template.html")
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

func (hd *CustomHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := hd.Repository.ListPosts(r.Context())
	if err != nil {
		handleError(w, "Failed to fetch posts", err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, posts, http.StatusOK)
}

func (hd *CustomHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var newPost model.NewPost

	if !decodeRequest(w, r, &newPost) {
		return
	}

	post, subscribers, newsletter, err := hd.Repository.CreatePost(r.Context(), newPost, r)

	if err != nil {
		handleError(w, "Failed to create post", err, http.StatusInternalServerError)
		return
	}

	cfg := utils.LoadConfig(".env")
	apiKey := cfg.ResendApiKey
	client := resend.NewClient(apiKey)
	// TODO: Change to production url
	unsubscribeUrl := "http://localhost:9069/api/unsubscribe"

	for _, subscriber := range subscribers {
		data := map[string]string{
			"title": post.Title,
			"content": post.Content,
			"unsubscribeUrl": unsubscribeUrl,
			"newsletterId":   newsletter.ID,
			"userId":         subscriber.ID,
		}
	
		mail, loadErr := loadEmailTemplate(data)
	
		if loadErr != nil {
			log.Printf("Failed load mail template: %v", err)
		}

        params := &resend.SendEmailRequest{
            From:    "newsletter@tapeer.cz",
            To:      []string{subscriber.Email},
            Subject: fmt.Sprintf("A new post in %s!", newsletter.Title),
            Html:    mail,
        }

        _, err := client.Emails.Send(params)
        if err != nil {
            log.Printf("Failed to send email to %s: %v", subscriber.Email, err)
        }
    }

	sendJSON(w, post, http.StatusOK)
}
