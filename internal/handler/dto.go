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

type linkStatsResponse struct {
	Code   string     `json:"code"`
	From   *time.Time `json:"from,omitempty"`
	To     *time.Time `json:"to,omitempty"`
	Clicks int64      `json:"clicks"`
}

type errorResponse struct {
	Error string `json:"error"`
}
