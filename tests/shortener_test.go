package shortener_test

import (
	"context"
	"errors"
	"testing"

	"url-shortener/internal/shortener"
	"url-shortener/internal/store"
)

type mockStore struct{}

func (mockStore) Create(context.Context, string, string) error { return nil }

func (mockStore) GetByCode(context.Context, string) (string, error) {
	return "", store.ErrNotFound
}

func TestCreate_ValidateURL(t *testing.T) {
	svc := shortener.New(mockStore{}, shortener.Options{BaseURL: "http://localhost:8080"})

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"https ok", "https://example.com/path", false},
		{"http ok", "http://example.com", false},
		{"no scheme", "example.com", true},
		{"ftp scheme", "ftp://example.com", true},
		{"empty", "", true},
		{"whitespace only", "   ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Create(context.Background(), tt.input)
			if tt.wantErr {
				if !errors.Is(err, shortener.ErrInvalidURL) {
					t.Errorf("Create(%q) error = %v, want ErrInvalidURL", tt.input, err)
				}
				return
			}
			if err != nil {
				t.Errorf("Create(%q) unexpected error: %v", tt.input, err)
			}
		})
	}
}

func TestCreate_CodeLength(t *testing.T) {
	svc := shortener.New(mockStore{}, shortener.Options{BaseURL: "http://localhost:8080"})

	result, err := svc.Create(context.Background(), "https://example.com")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	const wantLen = 7
	if len(result.Code) != wantLen {
		t.Errorf("len(code) = %d, want %d", len(result.Code), wantLen)
	}
}
