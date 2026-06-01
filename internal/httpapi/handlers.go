package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"url-shortener/internal/shortener"
	"url-shortener/internal/store"
)

type Handlers struct {
	svc *shortener.Service
}

func newHandlers(svc *shortener.Service) *Handlers {
	return &Handlers{svc: svc}
}

type CreateURLRequest struct {
	URL string `json:"url" example:"https://example.com/page"`
}

type CreateURLResponse struct {
	Code     string `json:"code" example:"aBc12Xy"`
	ShortURL string `json:"short_url" example:"http://localhost:8080/aBc12Xy"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"invalid url: must be http or https"`
}

// createURL godoc
// @Summary      Create a short URL
// @Description  Accepts a long HTTP/HTTPS URL and returns a short code and full short URL.
// @Tags         urls
// @Accept       json
// @Produce      json
// @Param        body  body      CreateURLRequest  true  "Long URL to shorten"
// @Success      201   {object}  CreateURLResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/urls [post]
func (h *Handlers) createURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid json body"})
		return
	}

	result, err := h.svc.Create(r.Context(), req.URL)
	if errors.Is(err, shortener.ErrInvalidURL) {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	writeJSON(w, http.StatusCreated, CreateURLResponse{
		Code:     result.Code,
		ShortURL: result.ShortURL,
	})
}

// redirect godoc
// @Summary      Resolve a short code
// @Description  Redirects to the original long URL for the given short code.
// @Tags         urls
// @Param        code  path  string  true  "Short code"
// @Success      302  {string}  string  "Redirect to long URL"
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /{code} [get]
func (h *Handlers) redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := r.PathValue("code")
	longURL, err := h.svc.Resolve(r.Context(), code)
	if errors.Is(err, store.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "not found"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}

// health godoc
// @Summary      Health check
// @Description  Returns ok when the service is running.
// @Tags         system
// @Produce      plain
// @Success      200  {string}  string  "ok"
// @Router       /health [get]
func (h *Handlers) health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
