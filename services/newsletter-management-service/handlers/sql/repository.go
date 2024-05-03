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
	err := pgxscan.Select(
		ctx,
		r.pool,
		&newsletters,
		ListNewslettersSQL,
	)
	return newsletters, err
}

func (r *Repository) CreateNewsletter(ctx context.Context, newsletter model.NewNewsletter, userId string) error {
	_, err := r.pool.Exec(
		ctx,
		CreateNewsletterSQL,
		newsletter.Title,
		newsletter.Description,
		userId,
	)

	return err
}

func (r *Repository) RenameNewsletter(ctx context.Context, newsletter model.UpdateNewsletter, userId string) error {
	_, err := r.pool.Exec(
		ctx,
		RenameNewsletterSQL,
		newsletter.Title,
		newsletter.Id,
		userId,
	)

	return err
}
