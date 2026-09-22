package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/mgreben/glink/internal/modules/links"
)

type linkService interface {
	GetByCode(ctx context.Context, code string) (*links.Link, error)
	Create(ctx context.Context, originalURL string) (*links.Link, error)
	RecordClick(ctx context.Context, linkID int64) error
	GetStats(ctx context.Context, code string, period links.StatsPeriod) (*links.LinkStats, error)
}

type LinkHandler struct {
	service  linkService
	validate *validator.Validate
}

func NewLinkHandler(validator *validator.Validate, service linkService) *LinkHandler {
	return &LinkHandler{
		service:  service,
		validate: validator,
	}
}

func (h *LinkHandler) RegisterRoutes(router chi.Router) {
	router.Route("/links", func(router chi.Router) {
		router.Post("/", h.Create)
		router.Get("/{code}/stats", h.Stats)
	})

	router.Get("/r/{code}", h.Redirect)
}

func (h *LinkHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request linkCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	request.normalize()
	if err := h.validate.Struct(request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid original_url")
		return
	}

	link, err := h.service.Create(r.Context(), request.OriginalURL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create link")
		return
	}

	writeJSON(w, http.StatusCreated, linkResponse{
		ID:          link.ID,
		OriginalURL: link.OriginalURL,
		Code:        link.Code,
		CreatedAt:   link.CreatedAt,
	})
}

func (h *LinkHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	link, err := h.service.GetByCode(r.Context(), code)
	if errors.Is(err, links.ErrNotFound) {
		writeError(w, http.StatusNotFound, "link not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get link")
		return
	}
	_ = h.service.RecordClick(r.Context(), link.ID)

	http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}

func (h *LinkHandler) Stats(w http.ResponseWriter, r *http.Request) {
	period, err := statsPeriod(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid from or to")
		return
	}

	stats, err := h.service.GetStats(r.Context(), chi.URLParam(r, "code"), period)
	if errors.Is(err, links.ErrNotFound) {
		writeError(w, http.StatusNotFound, "link not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get link stats")
		return
	}

	writeJSON(w, http.StatusOK, linkStatsResponse{
		Code:   stats.Code,
		From:   period.From,
		To:     period.To,
		Clicks: stats.Clicks,
	})
}

func statsPeriod(r *http.Request) (links.StatsPeriod, error) {
	from, err := parseOptionalTime(r.URL.Query().Get("from"))
	if err != nil {
		return links.StatsPeriod{}, err
	}
	to, err := parseOptionalTime(r.URL.Query().Get("to"))
	if err != nil {
		return links.StatsPeriod{}, err
	}
	if from != nil && to != nil && from.After(*to) {
		return links.StatsPeriod{}, errors.New("from is after to")
	}
	return links.StatsPeriod{From: from, To: to}, nil
}

func parseOptionalTime(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
