package sql

import (
	"context"
	"net/http"
	"publishing-service/handlers/model"
	"publishing-service/utils"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
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

func (r *Repository) CreatePost(ctx context.Context, post model.NewPost, req *http.Request) (model.Post, []model.Subscriber, model.Newsletter, error) {
	token := utils.GetBearerToken(req)
	userId, _ := utils.ExtractSubFromToken(token)

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
		return model.Post{}, []model.Subscriber{}, model.Newsletter{}, err
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
		return model.Post{}, []model.Subscriber{}, model.Newsletter{}, subscribersErr
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
		return model.Post{}, []model.Subscriber{}, model.Newsletter{}, newsletterErr
	}

	return createdPost, subscribers, newsletter, nil
}