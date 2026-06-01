package store

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("url not found")

type Store interface {
	Create(ctx context.Context, code, longURL string) error
	GetByCode(ctx context.Context, code string) (longURL string, err error)
}
