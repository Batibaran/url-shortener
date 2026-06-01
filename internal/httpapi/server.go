package httpapi

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	"url-shortener/internal/shortener"
)

func NewMux(svc *shortener.Service) *http.ServeMux {
	h := newHandlers(svc)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/urls", h.createURL)
	mux.HandleFunc("GET /health", h.health)
	mux.Handle("/swagger/", httpSwagger.Handler())
	mux.HandleFunc("GET /{code}", h.redirect)
	return mux
}
