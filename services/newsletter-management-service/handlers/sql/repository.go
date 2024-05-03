package sql

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"newsletter-management-service/handlers/model"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) ListNewsletters(ctx context.Context) ([]model.Newsletter, error) {
	var newsletters []model.Newsletter
	if err := pgxscan.Select(
		ctx,
		r.pool,
		&newsletters,
		ListNewslettersSQL,
	); err != nil {
		return nil, err
	}

	response := make([]model.Newsletter, len(newsletters))
	for i, newsletter := range newsletters {
		response[i] = model.Newsletter{
			ID:          newsletter.ID,
			CreatedAt:   newsletter.CreatedAt,
			Title:       newsletter.Title,
			Description: newsletter.Description,
			EditorID:    newsletter.EditorID,
		}
	}
	return response, nil
}

func (r *Repository) CreateNewsletter(ctx context.Context, newsletter model.NewNewsletter, userId string) error {
	if _, err := r.pool.Exec(
		ctx,
		CreateNewsletterSQL,
		newsletter.Title,
		newsletter.Description,
		userId,
	); err != nil {
		return err
	}

	return nil
}
