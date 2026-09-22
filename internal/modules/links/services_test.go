package links

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type linkRepoMock struct {
	mock.Mock
}

type linkCacheStub struct {
	link   *Link
	getErr error
	setErr error
}

func (c *linkCacheStub) Get(context.Context, string) (*Link, error) {
	return c.link, c.getErr
}

func (c *linkCacheStub) Set(context.Context, *Link) error {
	return c.setErr
}

func (l *linkRepoMock) GetByCode(ctx context.Context, code string) (*Link, error) {
	args := l.Called(ctx, code)
	if link, ok := args.Get(0).(*Link); ok {
		return link, args.Error(1)
	}
	return nil, args.Error(1)
}

func (l *linkRepoMock) Create(ctx context.Context, link *Link) error {
	args := l.Called(ctx, link)
	return args.Error(0)
}

func TestLinkService_GetByCode(t *testing.T) {
	t.Parallel()

	t.Run("returns link", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		expected := &Link{ID: 1, Code: "abc123"}
		repo := new(linkRepoMock)
		repo.On("GetByCode", ctx, "abc123").Return(expected, nil).Once()

		link, err := NewLinkService(repo, &linkCacheStub{}).GetByCode(ctx, "abc123")

		require.NoError(t, err)
		assert.Same(t, expected, link)
		repo.AssertExpectations(t)
	})

	t.Run("preserves domain error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		repo := new(linkRepoMock)
		repo.On("GetByCode", ctx, "missing").Return(nil, ErrNotFound).Once()

		link, err := NewLinkService(repo, &linkCacheStub{}).GetByCode(ctx, "missing")

		assert.Nil(t, link)
		assert.ErrorIs(t, err, ErrNotFound)
		repo.AssertExpectations(t)
	})
}

func TestLinkService_Create(t *testing.T) {
	t.Parallel()

	t.Run("creates link", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		repo := new(linkRepoMock)
		repo.On("Create", ctx, mock.MatchedBy(func(link *Link) bool {
			return link.OriginalURL == "https://example.com" && len(link.Code) == 6 && !link.CreatedAt.IsZero()
		})).Return(nil).Once()

		link, err := NewLinkService(repo, &linkCacheStub{}).Create(ctx, "https://example.com")

		require.NoError(t, err)
		assert.Equal(t, "https://example.com", link.OriginalURL)
		assert.Len(t, link.Code, 6)
		assert.WithinDuration(t, time.Now().UTC(), link.CreatedAt, time.Second)
		repo.AssertExpectations(t)
	})

	t.Run("retries a code collision", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		repo := new(linkRepoMock)
		repo.On("Create", ctx, mock.Anything).Return(ErrUniqueCode).Once()
		repo.On("Create", ctx, mock.Anything).Return(nil).Once()

		link, err := NewLinkService(repo, &linkCacheStub{}).Create(ctx, "https://example.com")

		require.NoError(t, err)
		assert.NotNil(t, link)
		repo.AssertExpectations(t)
	})

	t.Run("returns exhausted error after five collisions", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		repo := new(linkRepoMock)
		repo.On("Create", ctx, mock.Anything).Return(ErrUniqueCode).Times(5)

		link, err := NewLinkService(repo, &linkCacheStub{}).Create(ctx, "https://example.com")

		assert.Nil(t, link)
		assert.ErrorIs(t, err, ErrCodeGenerationExhausted)
		repo.AssertExpectations(t)
	})

	t.Run("wraps a repository error without retrying", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		repoErr := errors.New("database unavailable")
		repo := new(linkRepoMock)
		repo.On("Create", ctx, mock.Anything).Return(repoErr).Once()

		link, err := NewLinkService(repo, &linkCacheStub{}).Create(ctx, "https://example.com")

		assert.Nil(t, link)
		assert.ErrorIs(t, err, repoErr)
		repo.AssertExpectations(t)
	})
}

func TestLinkService_GetByCodeReturnsCachedLink(t *testing.T) {
	t.Parallel()
	expected := &Link{ID: 1, Code: "abc123", OriginalURL: "https://example.com"}

	link, err := NewLinkService(new(linkRepoMock), &linkCacheStub{link: expected}).GetByCode(context.Background(), "abc123")

	require.NoError(t, err)
	assert.Same(t, expected, link)
}

func TestLinkService_GetByCodeFallsBackWhenCacheFails(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	expected := &Link{ID: 1, Code: "abc123"}
	repo := new(linkRepoMock)
	repo.On("GetByCode", ctx, "abc123").Return(expected, nil).Once()

	link, err := NewLinkService(repo, &linkCacheStub{getErr: errors.New("redis unavailable")}).GetByCode(ctx, "abc123")

	require.NoError(t, err)
	assert.Same(t, expected, link)
	repo.AssertExpectations(t)
}
