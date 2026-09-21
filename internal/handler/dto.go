package handler

import (
	"strings"
	"time"
)

type linkCreateRequest struct {
	OriginalURL string `json:"original_url" validate:"required,http_url"`
}

func (r *linkCreateRequest) normalize() {
	r.OriginalURL = strings.TrimSpace(r.OriginalURL)
}

type linkResponse struct {
	ID          int64     `json:"id"`
	OriginalURL string    `json:"original_url"`
	Code        string    `json:"code"`
	CreatedAt   time.Time `json:"created_at"`
}

type errorResponse struct {
	Error string `json:"error"`
}
