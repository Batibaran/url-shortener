// @title           URL Shortener API
// @version         1.0
// @description     REST API for creating short links and redirecting to original URLs.
// @host            localhost:8080
// @BasePath        /
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	_ "url-shortener/docs"

	"url-shortener/internal/config"
	"url-shortener/internal/httpapi"
	"url-shortener/internal/shortener"
	"url-shortener/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database open: %v", err)
	}
	defer db.Close()

	if err := pingDB(db); err != nil {
		log.Fatalf("database ping: %v", err)
	}

	st := store.NewPostgres(db)
	svc := shortener.New(st, shortener.Options{BaseURL: cfg.BaseURL})
	mux := httpapi.NewMux(svc)

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func pingDB(db *sql.DB) error {
	const maxAttempts = 10
	const delay = time.Second

	var lastErr error
	for i := range maxAttempts {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		lastErr = db.PingContext(ctx)
		cancel()
		if lastErr == nil {
			return nil
		}
		if i < maxAttempts-1 {
			time.Sleep(delay)
		}
	}
	return fmt.Errorf("after %d attempts: %w", maxAttempts, lastErr)
}
