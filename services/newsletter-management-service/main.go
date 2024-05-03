package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	supa "github.com/nedpals/supabase-go"
	"log"
	"net/http"
	"newsletter-management-service/handlers"
	"newsletter-management-service/handlers/sql"
	"newsletter-management-service/utils"
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
		r.Get("/newsletters", hd.GetNewslettersHandler)
		r.Post("/newsletters", hd.CreateNewsletter)
		r.Put("/newsletters", hd.RenameNewsletter)
		r.Delete("/newsletters", hd.DeleteNewsletter)
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
