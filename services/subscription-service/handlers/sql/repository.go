package sql

import (
	"context"
	"subscription-service/handlers/model"

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

func (r *Repository) ListSubscriptions(ctx context.Context) ([]model.Subscription, error) {
	var subscriptions []model.Subscription
	if err := pgxscan.Select(
		ctx,
		r.pool,
		&subscriptions,
		ListSubscriptionsSql,
	); err != nil {
		return nil, err
	}

	response := make([]model.Subscription, len(subscriptions))
	for i, subscription := range subscriptions {
		response[i] = model.Subscription{
			ID: subscription.ID,
			CreatedAt: subscription.CreatedAt,
			NewsletterID: subscription.NewsletterID,
			SubscriberID: subscription.SubscriberID,
		}
	}
	return response, nil
}

func (r *Repository) Subscribe(ctx context.Context, newsletterId string, userId string) (model.Subscription, error) {
	var createdSubscription model.Subscription

	err := pgxscan.Get(
		ctx,
		r.pool,
		&createdSubscription,
		SubscribeSql,
		newsletterId,
		userId,
	)

	if err != nil {
		return model.Subscription{}, err
	}

	return createdSubscription, nil
}

func (r *Repository) Unsubcribe(ctx context.Context, newsletterId string, userId string) error {
	_, err := r.pool.Exec(
		ctx,
		UnsubcribeSql,
		newsletterId,
		userId,
	)

	return err
}