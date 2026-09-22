package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/mgreben/glink/internal/modules/links"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type linkServiceStub struct {
	stats       *links.LinkStats
	statsErr    error
	statsPeriod links.StatsPeriod
	statsCalls  int
}

func (s *linkServiceStub) GetByCode(context.Context, string) (*links.Link, error) { return nil, nil }
func (s *linkServiceStub) Create(context.Context, string) (*links.Link, error)    { return nil, nil }
func (s *linkServiceStub) RecordClick(context.Context, int64) error               { return nil }

func (s *linkServiceStub) GetStats(_ context.Context, _ string, period links.StatsPeriod) (*links.LinkStats, error) {
	s.statsCalls++
	s.statsPeriod = period
	return s.stats, s.statsErr
}

func TestLinkHandlerStats(t *testing.T) {
	from := time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, time.September, 30, 23, 59, 59, 0, time.UTC)
	service := &linkServiceStub{stats: &links.LinkStats{Code: "abc123", Clicks: 7}}
	router := chi.NewRouter()
	NewLinkHandler(validator.New(), service).RegisterRoutes(router)

	request := httptest.NewRequest(http.MethodGet, "/links/abc123/stats?from="+from.Format(time.RFC3339)+"&to="+to.Format(time.RFC3339), nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, 1, service.statsCalls)
	assert.Equal(t, &from, service.statsPeriod.From)
	assert.Equal(t, &to, service.statsPeriod.To)

	var body linkStatsResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
	assert.Equal(t, linkStatsResponse{Code: "abc123", From: &from, To: &to, Clicks: 7}, body)
}

func TestLinkHandlerStatsRejectsInvalidPeriod(t *testing.T) {
	service := &linkServiceStub{}
	router := chi.NewRouter()
	NewLinkHandler(validator.New(), service).RegisterRoutes(router)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/links/abc123/stats?from=not-a-date", nil))

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, 0, service.statsCalls)
}

func TestLinkHandlerStatsReturnsNotFound(t *testing.T) {
	service := &linkServiceStub{statsErr: fmt.Errorf("stats: %w", links.ErrNotFound)}
	router := chi.NewRouter()
	NewLinkHandler(validator.New(), service).RegisterRoutes(router)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/links/missing/stats", nil))

	assert.Equal(t, http.StatusNotFound, response.Code)
}
