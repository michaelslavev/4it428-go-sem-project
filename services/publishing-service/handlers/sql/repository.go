package sql

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"publishing-service/handlers/model"
	"publishing-service/utils"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/resend/resend-go/v2"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) ListPosts(ctx context.Context) ([]model.Post, error) {
	var posts []model.Post
	if err := pgxscan.Select(
		ctx,
		r.pool,
		&posts,
		ListPostsSql,
	); err != nil {
		return nil, err
	}

	response := make([]model.Post, len(posts))
	for i, post := range posts {
		response[i] = model.Post{
			ID: post.ID,
			CreatedAt: post.CreatedAt,
			Title: post.Title,
			Content: post.Content,
			NewsletterID: post.NewsletterID,
		}
	}
	return response, nil
}

func (r *Repository) CreatePost(ctx context.Context, post model.NewPost, req *http.Request) (model.Post, error) {
	token := utils.GetBearerToken(req)
	userId, _ := utils.ExtractSubFromToken(token)
	cfg := utils.LoadConfig(".env")

	var createdPost model.Post

	err := pgxscan.Get(
		ctx,
		r.pool,
		&createdPost,
		CreatePostSql,
		post.Title,
		post.Content,
		post.NewsletterID,
	);

	if err != nil {
		return model.Post{}, err
	}

	var subscribers []model.Subscriber

	subscribersErr := pgxscan.Select(
		ctx,
		r.pool,
		&subscribers,
		GetNewsletterSubscribersSql,
		post.NewsletterID,
		userId,
	)

	if subscribersErr != nil {
		return model.Post{}, subscribersErr
	}

	var newsletter model.Newsletter

	newsletterErr := pgxscan.Get(
		ctx,
		r.pool,
		&newsletter,
		GetNewsletter,
		post.NewsletterID,
	)

	if newsletterErr != nil {
		return model.Post{}, newsletterErr
	}

	apiKey := cfg.ResendApiKey
	client := resend.NewClient(apiKey)

	for _, subscriber := range subscribers {
        params := &resend.SendEmailRequest{
            From:    "newsletter@tapeer.cz",
            To:      []string{subscriber.Email},
            Subject: newsletter.Title,
            Html:    fmt.Sprintf("<div><strong>New post available: %s</strong><br /><p>%s</p></div>", post.Title, post.Content),
        }

        _, err := client.Emails.Send(params)
        if err != nil {
            log.Printf("Failed to send email to %s: %v", subscriber.Email, err)
        }
    }

	return createdPost, nil
}