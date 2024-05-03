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

func (r *Repository) CreateNewsletter(ctx context.Context, newsletter model.NewNewsletter, userId string) (model.Newsletter, error) {
	var createdNewsletter model.Newsletter
	err := pgxscan.Get(
		ctx,
		r.pool,
		&createdNewsletter,
		CreateNewsletterSQL,
		newsletter.Title,
		newsletter.Description,
		userId,
	)
	if err != nil {
		return model.Newsletter{}, err
	}
	return createdNewsletter, nil
}

func (r *Repository) RenameNewsletter(ctx context.Context, newsletter model.UpdateNewsletter, userId string) (model.Newsletter, error) {
	var updatedNewsLetter model.Newsletter
	err := pgxscan.Get(
		ctx,
		r.pool,
		&updatedNewsLetter,
		RenameNewsletterSQL,
		newsletter.Title,
		newsletter.Id,
		userId,
	)
	if err != nil {
		return model.Newsletter{}, err
	}
	return updatedNewsLetter, nil
}
