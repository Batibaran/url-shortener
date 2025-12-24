package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Storage holds the database connection pool.
type Storage struct {
	db *pgxpool.Pool
}

// URL represents a record in the urls table.
type URL struct {
	ID          int64  `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortCode   string `json:"short_code"`
}

// New creates a new Storage and connects to the database.
func New(databaseURL string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

// SaveURL saves a new URL to the database.
func (s *Storage) SaveURL(originalURL, shortCode string) error {
	query := "INSERT INTO urls (original_url, short_code) VALUES ($1, $2)"
	_, err := s.db.Exec(context.Background(), query, originalURL, shortCode)
	return err
}

// GetURLByShortCode retrieves the original URL for a given short code.
func (s *Storage) GetURLByShortCode(shortCode string) (string, error) {
	var originalURL string
	query := "SELECT original_url FROM urls WHERE short_code = $1"
	err := s.db.QueryRow(context.Background(), query, shortCode).Scan(&originalURL)
	if err != nil {
		log.Printf("Error retrieving URL for short code %s: %v", shortCode, err)
		return "", err
	}
	return originalURL, nil
}

// Close closes the database connection.
func (s *Storage) Close() {
	s.db.Close()
}
