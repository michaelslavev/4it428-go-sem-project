package sql

import (
	"context"
	"publishing-service/handlers/model"

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

func (r *Repository) CreatePost(ctx context.Context, post model.NewPost) error {
	if _, err := r.pool.Exec(
		ctx,
		CreatePostSql,
		post.Title,
		post.Content,
		post.NewsletterID,
	); err != nil {
		return err
	}

	return nil
}
