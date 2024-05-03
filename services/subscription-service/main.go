package main

import (
	"context"
	"log"
	"net/http"
	"subscription-service/handlers"
	"subscription-service/handlers/sql"
	"subscription-service/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	supa "github.com/nedpals/supabase-go"
)

var ctx context.Context
var cfg utils.ServerConfig
var supabase *supa.Client
var database *pgxpool.Pool
var repository *sql.Repository

func init() {
	ctx = context.Background()
	cfg = utils.LoadConfig(".env")
	supabase = supa.CreateClient(cfg.SupabaseURL, cfg.SupabaseKEY)
	database, _ = setupDatabase(ctx, cfg)
	repository = sql.NewRepository(database)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	hd := handlers.NewCustomHandler(supabase, repository)
	r.Route("/api", func(r chi.Router) {
		r.Get("/subscriptions", hd.GetSubscriptions)
		r.Post("/subscribe/{id}", hd.Subscribe)
		r.Post("/unsubscribe/{id}", hd.Unsubcribe)
	})

	// Starting server
	address := cfg.IP + ":" + cfg.Port
	log.Printf("Server starting on %s", address)
	err := http.ListenAndServe(address, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func setupDatabase(ctx context.Context, cfg utils.ServerConfig) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(
		ctx,
		cfg.DatabaseURL,
	)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
