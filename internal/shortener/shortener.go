package shortener

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"strings"

	"url-shortener/internal/store"
)

const (
	codeLength = 7
	codeChars  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	maxRetries = 5
)

type Service struct {
	store   store.Store
	baseURL string
}

type Options struct {
	BaseURL string
}

func New(st store.Store, opts Options) *Service {
	return &Service{
		store:   st,
		baseURL: strings.TrimRight(opts.BaseURL, "/"),
	}
}

type CreateResult struct {
	Code     string
	ShortURL string
}

func (s *Service) Create(ctx context.Context, rawURL string) (CreateResult, error) {
	longURL, err := validateURL(rawURL)
	if err != nil {
		return CreateResult{}, err
	}

	for range maxRetries {
		code, err := generateCode(codeLength)
		if err != nil {
			return CreateResult{}, fmt.Errorf("generate code: %w", err)
		}

		err = s.store.Create(ctx, code, longURL)
		if err == nil {
			return CreateResult{
				Code:     code,
				ShortURL: s.baseURL + "/" + code,
			}, nil
		}
		if !errors.Is(err, store.ErrDuplicateCode) {
			return CreateResult{}, err
		}
	}

	return CreateResult{}, errors.New("could not generate unique code")
}

func (s *Service) Resolve(ctx context.Context, code string) (string, error) {
	if code == "" {
		return "", store.ErrNotFound
	}
	return s.store.GetByCode(ctx, code)
}

var ErrInvalidURL = errors.New("invalid url: must be http or https")

func validateURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	u, err := url.Parse(raw)
	if err != nil {
		return "", ErrInvalidURL
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", ErrInvalidURL
	}
	if u.Host == "" {
		return "", ErrInvalidURL
	}
	return u.String(), nil
}

func generateCode(length int) (string, error) {
	b := make([]byte, length)
	max := big.NewInt(int64(len(codeChars)))
	for i := range b {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = codeChars[n.Int64()]
	}
	return string(b), nil
}
