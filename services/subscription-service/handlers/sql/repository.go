package sql

import (
	"context"
	"fmt"
	"log"
	"subscription-service/handlers/model"
	"subscription-service/utils"

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

func (r *Repository) Subscribe(ctx context.Context, newsletterId string, userId string) (model.Subscription, error) {
	// Check if a subscription already exists
    var count int

    checkErr := pgxscan.Get(ctx, r.pool, &count, CheckSubscriptionsSql, newsletterId, userId)

    if checkErr != nil {
        return model.Subscription{}, fmt.Errorf("error checking existing subscription: %v", checkErr)
    }

    // If the subscription count is greater than 0, return an error indicating that the subscription already exists
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

	var subscriber model.Subscriber

	subscriberErr := pgxscan.Get(
		ctx,
		r.pool,
		&subscriber,
		GetUserByIdSql,
		userId,
	)

	if subscriberErr != nil {
		return model.Subscription{}, subscriberErr
	}

	cfg := utils.LoadConfig(".env")
	apiKey := cfg.ResendApiKey
	client := resend.NewClient(apiKey)
	unsubscribeUrl := "http://localhost:9069/api/unsubscribe"

	params := &resend.SendEmailRequest{
		From:    "newsletter@tapeer.cz",
		To:      []string{subscriber.Email},
		Subject: "You have been subscribed.",
		Html:    fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome to the club!</h1>
				<p>We are excited to have you onboard and look forward to bringing you the latest updates directly to your inbox.</p>
				<p>If you have any questions, feel free to contact us anytime.</p>
				<hr>
				<footer>
					<p>If you wish to unsubscribe and stop receiving these emails, please click on the link below:</p>
					<p><a href="%s?newsletterId=%s&userId=%s" style="color: #1155cc;">Unsubscribe</a></p>
				</footer>
			</body>
		</html>`, unsubscribeUrl, newsletterId, userId),
		Headers: map[string]string{
			"List-Unsubscribe": fmt.Sprintf("<%s?newsletterId=%s&userId=%s>", unsubscribeUrl, newsletterId, userId),
			"List-Unsubscribe-Post": "List-Unsubscribe=One-Click",
		},
	}

	_, emailErr := client.Emails.Send(params)
	if emailErr != nil {
		log.Printf("Failed to send email to %s: %v", subscriber.Email, err)
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