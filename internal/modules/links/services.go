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

type LinkService struct {
	repo linkRepo
}

func NewLinkService(repo linkRepo) *LinkService {
	return &LinkService{
		repo: repo,
	}
}

func (s *LinkService) GetByCode(ctx context.Context, code string) (*Link, error) {
	link, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get link by code: %w", err)
	}
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
			return link, nil
		}

		if !errors.Is(err, ErrUniqueCode) {
			return nil, fmt.Errorf("failed to save link: %w", err)
		}
	}

	return nil, fmt.Errorf("create link: %w", ErrCodeGenerationExhausted)
}
