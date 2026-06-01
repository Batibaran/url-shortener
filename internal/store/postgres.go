package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

func (p *Postgres) Create(ctx context.Context, code, longURL string) error {
	_, err := p.db.ExecContext(ctx,
		`INSERT INTO urls (code, long_url) VALUES ($1, $2)`,
		code, longURL,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicateCode
		}
		return err
	}
	return nil
}

var ErrDuplicateCode = errors.New("duplicate code")

func (p *Postgres) GetByCode(ctx context.Context, code string) (string, error) {
	var longURL string
	err := p.db.QueryRowContext(ctx,
		`SELECT long_url FROM urls WHERE code = $1`,
		code,
	).Scan(&longURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return longURL, nil
}
