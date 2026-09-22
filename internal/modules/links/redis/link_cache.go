package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mgreben/glink/internal/config"
	"github.com/mgreben/glink/internal/modules/links"
	redis "github.com/redis/go-redis/v9"
)

const linkKeyPrefix = "link:"

type LinkCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewLinkCache(client *redis.Client, cfg *config.Config) *LinkCache {
	return &LinkCache{client: client, ttl: cfg.Redis.LinkCacheTTL}
}

func (c *LinkCache) Get(ctx context.Context, code string) (*links.Link, error) {
	value, err := c.client.Get(ctx, linkKey(code)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get cached link: %w", err)
	}

	link := &links.Link{}
	if err := json.Unmarshal(value, link); err != nil {
		return nil, fmt.Errorf("decode cached link: %w", err)
	}
	return link, nil
}

func (c *LinkCache) Set(ctx context.Context, link *links.Link) error {
	value, err := json.Marshal(link)
	if err != nil {
		return fmt.Errorf("encode cached link: %w", err)
	}
	if err := c.client.Set(ctx, linkKey(link.Code), value, c.ttl).Err(); err != nil {
		return fmt.Errorf("cache link: %w", err)
	}
	return nil
}

func linkKey(code string) string {
	return linkKeyPrefix + code
}
