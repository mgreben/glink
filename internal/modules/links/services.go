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
	RecordClick(context.Context, int64, time.Time) error
	GetStats(context.Context, string, StatsPeriod) (*LinkStats, error)
}

func (s *LinkService) RecordClick(ctx context.Context, linkID int64) error {
	if err := s.repo.RecordClick(ctx, linkID, time.Now().UTC()); err != nil {
		return fmt.Errorf("record link click: %w", err)
	}
	return nil
}

func (s *LinkService) GetStats(ctx context.Context, code string, period StatsPeriod) (*LinkStats, error) {
	stats, err := s.repo.GetStats(ctx, code, period)
	if err != nil {
		return nil, fmt.Errorf("get link stats: %w", err)
	}
	return stats, nil
}

type linkCache interface {
	Get(context.Context, string) (*Link, error)
	Set(context.Context, *Link) error
}

type LinkService struct {
	repo  linkRepo
	cache linkCache
}

func NewLinkService(repo linkRepo, cache linkCache) *LinkService {
	return &LinkService{
		repo:  repo,
		cache: cache,
	}
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
