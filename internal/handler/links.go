package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/mgreben/glink/internal/modules/links"
)

type linkService interface {
	GetByCode(ctx context.Context, code string) (*links.Link, error)
	Create(ctx context.Context, originalURL string) (*links.Link, error)
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

	http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}
