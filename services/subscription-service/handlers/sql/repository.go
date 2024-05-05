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

func (r *Repository) Subscribe(ctx context.Context, newsletterId string, userId string) (model.Subscription, model.Subscriber, model.Newsletter, error) {
	// Check if a subscription already exists
    var count int

    checkErr := pgxscan.Get(ctx, r.pool, &count, CheckSubscriptionsSql, newsletterId, userId)

    if checkErr != nil {
        return model.Subscription{}, model.Subscriber{}, model.Newsletter{}, fmt.Errorf("error checking existing subscription: %v", checkErr)
    }

    // If the subscription count is greater than 0, return an error indicating that the subscription already exists
    if count > 0 {
        return model.Subscription{}, model.Subscriber{}, model.Newsletter{}, utils.ErrSubscriptionExists
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
		return model.Subscription{}, model.Subscriber{}, model.Newsletter{}, err
	}

	// Select a subscriber to get his email
	var subscriber model.Subscriber

	subscriberErr := pgxscan.Get(
		ctx,
		r.pool,
		&subscriber,
		GetUserByIdSql,
		userId,
	)

	if subscriberErr != nil {
		return model.Subscription{}, model.Subscriber{}, model.Newsletter{}, subscriberErr
	}

	// Select a newsletter to get its title
	var newsletter model.Newsletter

	newsletterErr := pgxscan.Get(
		ctx,
		r.pool,
		&newsletter,
		GetNewsletterById,
		newsletterId,
	)

	if newsletterErr != nil {
		return model.Subscription{}, model.Subscriber{}, model.Newsletter{}, newsletterErr
	}

	return createdSubscription, subscriber, newsletter, nil
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