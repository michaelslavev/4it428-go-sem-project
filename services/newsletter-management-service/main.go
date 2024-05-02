package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	supa "github.com/nedpals/supabase-go"
	"log"
	"log/slog"
	"net/http"
	"newsletter-management-service/handlers"
	"newsletter-management-service/handlers/sql"
	"newsletter-management-service/utils"
)

func main() {
	ctx := context.Background()
	cfg := utils.LoadConfig(".env")
	supabase := supa.CreateClient(cfg.SupabaseURL, cfg.SupabaseKEY)

	database, err := setupDatabase(ctx, cfg)
	if err != nil {
		slog.Error("initializing database", slog.Any("error", err))
	}
	repository := sql.NewRepository(database)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	hd := handlers.NewCustomHandler(supabase, repository)
	r.Route("/api", func(r chi.Router) {
		r.Get("/listNewsletters", hd.GetNewslettersHandler)
		r.Post("/createNewsletter", hd.CreateNewsletter)
		r.Post("/renameNewsletter", hd.RenameNewsletter)
		r.Post("/deleteNewsletter", hd.DeleteNewsletter)
	})

	// Starting server
	address := cfg.IP + ":" + cfg.Port
	log.Printf("Server starting on %s", address)
	err = http.ListenAndServe(address, r)
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
