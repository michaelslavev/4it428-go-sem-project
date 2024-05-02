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
		ListNewsletters,
	); err != nil {
		return nil, err
	}

	response := make([]model.Newsletter, len(newsletters))
	for i, newsletter := range newsletters {
		response[i] = model.Newsletter{
			ID: newsletter.ID,
		}
	}
	return response, nil
}
