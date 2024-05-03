package sql

import (
	"context"
	"fmt"
	"subscription-service/handlers/model"
	"subscription-service/utils"

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

func (r *Repository) Subscribe(ctx context.Context, newsletterId string, userId string) (model.Subscription, error) {
	// Check if a subscription already exists
    var count int

    checkErr := pgxscan.Get(ctx, r.pool, &count, CheckSubscriptionsSql, newsletterId, userId)

    if checkErr != nil {
        return model.Subscription{}, fmt.Errorf("error checking existing subscription: %v", checkErr)
    }

    // If the subscription count is greater than 0, return an error indicating the subscription already exists
    if count > 0 {
        return model.Subscription{}, utils.ErrSubscriptionExists
    }

    // If no existing subscription, proceed to create a new one
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