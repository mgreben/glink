package links

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mgreben/glink/pkg/base62"
)

type linkRepo interface {
	GetByCode(context.Context, string) (*Link, error)
	Create(context.Context, *Link) error
}

type statsReader interface {
	CountClicks(context.Context, int64, StatsPeriod) (int64, error)
}

type clickPublisher interface {
	Publish(context.Context, ClickEvent)
}

func (s *LinkService) RecordClick(ctx context.Context, linkID int64) error {
	s.publisher.Publish(ctx, ClickEvent{LinkID: linkID, ClickedAt: time.Now().UTC()})
	return nil
}

func (s *LinkService) GetStats(ctx context.Context, code string, period StatsPeriod) (*LinkStats, error) {
	link, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("get link stats: %w", err)
	}
	clicks, err := s.stats.CountClicks(ctx, link.ID, period)
	if err != nil {
		return nil, fmt.Errorf("get link stats: %w", err)
	}
	return &LinkStats{Code: link.Code, Clicks: clicks}, nil
}

type linkCache interface {
	Get(context.Context, string) (*Link, error)
	Set(context.Context, *Link) error
}

type LinkService struct {
	repo      linkRepo
	cache     linkCache
	publisher clickPublisher
	stats     statsReader
}

func NewLinkService(repo linkRepo, cache linkCache, publisher clickPublisher) *LinkService {
	return &LinkService{
		repo:      repo,
		cache:     cache,
		publisher: publisher,
	}
}

func NewLinkServiceWithStats(repo linkRepo, cache linkCache, publisher clickPublisher, stats statsReader) *LinkService {
	service := NewLinkService(repo, cache, publisher)
	service.stats = stats
	return service
}

func (s *LinkService) GetByCode(ctx context.Context, code string) (*Link, error) {
	link, err := s.cache.Get(ctx, code)
	if err == nil && link != nil {
		return link, nil
	}

	link, err = s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get link by code: %w", err)
	}
	_ = s.cache.Set(ctx, link)
	return link, nil
}

func (s *LinkService) Create(ctx context.Context, originalURL string) (*Link, error) {
	link := &Link{
		OriginalURL: originalURL,
		CreatedAt:   time.Now().UTC(),
	}

	for attempt := 0; attempt < 5; attempt++ {
		code, err := base62.GenerateCode(6)
		if err != nil {
			return nil, fmt.Errorf("failed to generate link code: %w", err)
		}

		link.Code = code

		err = s.repo.Create(ctx, link)
		if err == nil {
			_ = s.cache.Set(ctx, link)
			return link, nil
		}

		if !errors.Is(err, ErrUniqueCode) {
			return nil, fmt.Errorf("failed to save link: %w", err)
		}
	}

	return nil, fmt.Errorf("create link: %w", ErrCodeGenerationExhausted)
}
